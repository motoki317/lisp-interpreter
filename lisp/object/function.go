package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewWrappedFunctionObject(f func(objects []Object) Object) Object {
	return function(func(objects []Object) (Object, *node.Node, *Env) {
		return f(objects), nil, nil
	})
}

func NewFunctionObject(f func(objects []Object) (Object, *node.Node, *Env)) Object {
	return function(f)
}

func (f function) Type() object_type.T {
	return object_type.Function
}

func (f function) Number() float64 {
	panic("number() called on function object")
}

func (f function) Bool() bool {
	panic("Bool() called on function object")
}

func (f function) Pair() *[2]Object {
	panic("Pair() called on function object")
}

func (f function) Str() string {
	panic("Str() called on function object")
}

func (f function) F(objects []Object) (Object, *node.Node, *Env) {
	return f(objects)
}

func (f function) String() string {
	return "<function>"
}

func (f function) Display() string {
	return "<function>"
}

func (f function) IsList() bool {
	return false
}

func (f function) ListElements() []Object {
	panic("ListElements() called on function object")
}

func (f function) IsTruthy() bool {
	return true
}

func (f function) Equals(object Object) bool {
	if object.Type() != object_type.Function {
		return false
	}
	// function comparison undefined
	return true
}
