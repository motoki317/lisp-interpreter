package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func (v void) Type() object_type.T {
	return object_type.Void
}

func (v void) Number() float64 {
	panic("number() called on void object")
}

func (v void) Bool() bool {
	panic("Bool() called on void object")
}

func (v void) Pair() *[2]Object {
	panic("Pair() called on void object")
}

func (v void) Str() string {
	panic("Str() called on void object")
}

func (v void) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on void object")
}

func (v void) String() string {
	return "<void>"
}

func (v void) Display() string {
	return ""
}

func (v void) IsList() bool {
	return false
}

func (v void) ListElements() []Object {
	panic("ListElements() called on void object")
}

func (v void) IsTruthy() bool {
	return true
}

func (v void) Equals(object Object) bool {
	return object.Type() == object_type.Void
}
