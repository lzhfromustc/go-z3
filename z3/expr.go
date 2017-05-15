// Copyright 2017 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package z3

import (
	"math/big"
	"runtime"
	"unsafe"
)

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// TODO: Should the various Value types disallow explicit conversion?
// Right now they all have the same underlying type.

// An Value is a symbolic value (a Z3 expression).
//
// This package exports a concrete type for each different kind of
// value, such as Bool, BV, and Int. These concrete types provide
// methods for deriving other values.
//
// Having separate types for each kind separates which methods can be
// applied to which kind of value and provides some level of static
// type safety. However, by no means does this fully capture Z3's type
// system, so dynamic type checking can still fail.
type Value interface {
	// AsAST returns the abstract syntax tree underlying this
	// Value.
	AsAST() AST

	// Sort returns this value's sort.
	Sort() Sort

	// String returns an S-expression representation of this value.
	String() string

	astKind() C.Z3_ast_kind
	impl() *exprImpl
}

type noEq struct {
	_ [0]func()
}

// expr is a general wrapper for the Z3_ast type. Expression values
// are implemented as public types corresponding to Z3 sorts that are
// named types for expr.
type expr struct {
	// *exprImpl is the internal state of expr. This is wrapped
	// and unexported so we can attach a finalizer to the exprImpl
	// object without any possibility of user code copying the
	// underlying wrapper and breaking our tracking.
	*exprImpl

	// noEq prevents user code from directly comparing exprs for
	// equality.
	noEq
}

type exprImpl astImpl

func wrapExpr(ctx *Context, c C.Z3_ast) expr {
	return expr{(*exprImpl)(wrapAST(ctx, c).astImpl), noEq{}}
}

// lift wraps x in the appropriate Value type. kind must be x's kind if
// known or otherwise SortUnknown.
func (x expr) lift(kind Kind) Value {
	if kind == KindUnknown {
		kind = x.Sort().Kind()
	}
	wrap, ok := kindWrappers[kind]
	if !ok {
		panic("expression has unknown kind " + kind.String())
	}
	return wrap(x)
}

// Const returns a constant named "name" of the given sort. This
// constant will be same as all other constants created with this
// name.
func (ctx *Context) Const(name string, sort Sort) Value {
	sym := ctx.symbol(name)
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_const(ctx.c, sym, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr).lift(sort.Kind())
}

// FreshConst returns a constant that is distinct from all other
// constants. The name will begin with "prefix".
func (ctx *Context) FreshConst(prefix string, sort Sort) Value {
	cprefix := C.CString(prefix)
	defer C.free(unsafe.Pointer(cprefix))
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_fresh_const(ctx.c, cprefix, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr).lift(sort.Kind())
}

// FromBigInt returns a literal whose value is val. sort must have
// kind int, real, finite-domain, or bit-vector.
func (ctx *Context) FromBigInt(val *big.Int, sort Sort) Value {
	cstr := C.CString(val.Text(10))
	defer C.free(unsafe.Pointer(cstr))
	var cexpr C.Z3_ast
	ctx.do(func() {
		cexpr = C.Z3_mk_numeral(ctx.c, cstr, sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr).lift(sort.Kind())
}

// FromInt returns a literal whose value is val. sort must have kind
// int, real, finite-domain, or bit-vector.
func (ctx *Context) FromInt(val int64, sort Sort) Value {
	var cexpr C.Z3_ast
	ctx.do(func() {
		// Z3_mk_int64 doesn't say real sorts are accepted,
		// but the C++ bindings use it for reals.
		cexpr = C.Z3_mk_int64(ctx.c, C.__int64(val), sort.c)
	})
	runtime.KeepAlive(sort)
	return wrapExpr(ctx, cexpr).lift(sort.Kind())
}

func (expr *exprImpl) impl() *exprImpl {
	return expr
}

// String returns a string representation of expr.
func (expr *exprImpl) String() string {
	var res string
	expr.ctx.do(func() {
		res = C.GoString(C.Z3_ast_to_string(expr.ctx.c, expr.c))
	})
	runtime.KeepAlive(expr)
	return res
}

// AsAST returns the abstract syntax tree underlying expr.
func (expr *exprImpl) AsAST() AST {
	ast := AST{(*astImpl)(expr), noEq{}}
	runtime.KeepAlive(expr)
	return ast
}

// Sort returns expr's sort.
func (expr *exprImpl) Sort() Sort {
	var csort C.Z3_sort
	expr.ctx.do(func() {
		csort = C.Z3_get_sort(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return wrapSort(expr.ctx, csort, KindUnknown)
}

func (expr *exprImpl) astKind() C.Z3_ast_kind {
	var ckind C.Z3_ast_kind
	expr.ctx.do(func() {
		ckind = C.Z3_get_ast_kind(expr.ctx.c, expr.c)
	})
	runtime.KeepAlive(expr)
	return ckind
}

func (expr *exprImpl) asBigInt() (val *big.Int, isLiteral bool) {
	switch expr.Sort().Kind() {
	default:
		panic("sort " + expr.Sort().String() + " cannot be represented as a big.Int")
	case KindInt, KindBV:
	}
	if expr.astKind() != C.Z3_NUMERAL_AST {
		return nil, false
	}
	var str string
	expr.ctx.do(func() {
		cstr := C.Z3_get_numeral_string(expr.ctx.c, expr.c)
		str = C.GoString(cstr)
	})
	var v big.Int
	if _, ok := v.SetString(str, 10); !ok {
		panic("failed to parse numeral string")
	}
	return &v, true
}

func (expr *exprImpl) asInt64() (val int64, isLiteral, ok bool) {
	switch expr.Sort().Kind() {
	default:
		panic("sort " + expr.Sort().String() + " cannot be represented as an int64")
	case KindInt, KindBV:
	}
	if expr.astKind() != C.Z3_NUMERAL_AST {
		return 0, false, false
	}
	var cval C.__int64
	expr.ctx.do(func() {
		ok = z3ToBool(C.Z3_get_numeral_int64(expr.ctx.c, expr.c, &cval))
	})
	return int64(cval), true, ok
}

func (expr *exprImpl) asUint64() (val uint64, isLiteral, ok bool) {
	switch expr.Sort().Kind() {
	default:
		panic("sort " + expr.Sort().String() + " cannot be represented as an int64")
	case KindInt, KindBV:
	}
	if expr.astKind() != C.Z3_NUMERAL_AST {
		return 0, false, false
	}
	var cval C.__uint64
	expr.ctx.do(func() {
		ok = z3ToBool(C.Z3_get_numeral_uint64(expr.ctx.c, expr.c, &cval))
	})
	return uint64(cval), true, ok
}
