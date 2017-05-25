// Generated by genwrap.go. DO NOT EDIT

package z3

import "runtime"

/*
#cgo LDFLAGS: -lz3
#include <z3.h>
#include <stdlib.h>
*/
import "C"

// Eq returns a Value that is true if l and r are equal.
func (l Real) Eq(r Real) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l Real) NE(r Real) Bool {
	return l.ctx.Distinct(l, r)
}

// Div returns l / r.
//
// If r is 0, the result is unconstrained.
func (l Real) Div(r Real) Real {
	// Generated from real.go:124.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_div(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Real(val)
}

// ToInt returns the floor of l as sort Int.
//
// Note that this is not truncation. For example, ToInt(-1.3) is -2.
func (l Real) ToInt() Int {
	// Generated from real.go:130.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_real2int(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Int(val)
}

// IsInt returns a Value that is true if l has no fractional part.
func (l Real) IsInt() Bool {
	// Generated from real.go:134.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_is_int(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Bool(val)
}

// Add returns the sum l + r[0] + r[1] + ...
func (l Real) Add(r ...Real) Real {
	// Generated from intreal.go:12.
	ctx := l.ctx
	cargs := make([]C.Z3_ast, len(r)+1)
	cargs[0] = l.c
	for i, arg := range r {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_add(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return Real(val)
}

// Mul returns the product l * r[0] * r[1] * ...
func (l Real) Mul(r ...Real) Real {
	// Generated from intreal.go:16.
	ctx := l.ctx
	cargs := make([]C.Z3_ast, len(r)+1)
	cargs[0] = l.c
	for i, arg := range r {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_mul(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return Real(val)
}

// Sub returns l - r[0] - r[1] - ...
func (l Real) Sub(r ...Real) Real {
	// Generated from intreal.go:20.
	ctx := l.ctx
	cargs := make([]C.Z3_ast, len(r)+1)
	cargs[0] = l.c
	for i, arg := range r {
		cargs[i+1] = arg.c
	}
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_sub(ctx.c, C.uint(len(cargs)), &cargs[0])
	})
	runtime.KeepAlive(&cargs[0])
	return Real(val)
}

// Neg returns -l.
func (l Real) Neg() Real {
	// Generated from intreal.go:24.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_unary_minus(ctx.c, l.c)
	})
	runtime.KeepAlive(l)
	return Real(val)
}

// Exp returns lᶠ.
func (l Real) Exp(r Real) Real {
	// Generated from intreal.go:28.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_power(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Real(val)
}

// LT returns l < r.
func (l Real) LT(r Real) Bool {
	// Generated from intreal.go:32.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_lt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// LE returns l <= r.
func (l Real) LE(r Real) Bool {
	// Generated from intreal.go:36.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_le(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// GT returns l > r.
func (l Real) GT(r Real) Bool {
	// Generated from intreal.go:40.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_gt(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// GE returns l >= r.
func (l Real) GE(r Real) Bool {
	// Generated from intreal.go:44.
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_ge(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}
