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
func (l Uninterpreted) Eq(r Uninterpreted) Bool {
	ctx := l.ctx
	val := wrapValue(ctx, func() C.Z3_ast {
		return C.Z3_mk_eq(ctx.c, l.c, r.c)
	})
	runtime.KeepAlive(l)
	runtime.KeepAlive(r)
	return Bool(val)
}

// NE returns a Value that is true if l and r are not equal.
func (l Uninterpreted) NE(r Uninterpreted) Bool {
	return l.ctx.Distinct(l, r)
}
