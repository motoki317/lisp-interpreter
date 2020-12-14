package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func (n null) Type() object_type.T {
	return object_type.Null
}

func (n null) Number() float64 {
	panic("number() called on null object")
}

func (n null) Bool() bool {
	panic("Bool() called on null object")
}

func (n null) Pair() *[2]Object {
	panic("Pair() called on null object")
}

func (n null) Str() string {
	panic("Str() called on null object")
}

func (n null) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on null object")
}

func (n null) String() string {
	return "()"
}

func (n null) Display() string {
	return "()"
}

func (n null) IsList() bool {
	// a null object is a list
	return true
}

func (n null) ListElements() []Object {
	// base case, a null object is an empty list
	return []Object{}
}

func (n null) IsTruthy() bool {
	return true
}

func (n null) Equals(object Object) bool {
	return object.Type() == object_type.Null
}
