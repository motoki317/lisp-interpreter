package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewStringObject(s string) Object {
	return str(s)
}

func (s str) Type() object_type.T {
	return object_type.Str
}

func (s str) Number() float64 {
	panic("number() called on str object")
}

func (s str) Bool() bool {
	panic("Bool() called on str object")
}

func (s str) Pair() *[2]Object {
	panic("Pair() called on str object")
}

func (s str) Str() string {
	return string(s)
}

func (s str) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on str object")
}

func (s str) String() string {
	return "\"" + string(s) + "\""
}

func (s str) IsList() bool {
	return false
}

func (s str) ListElements() []Object {
	panic("ListElements() called on str object")
}

func (s str) IsTruthy() bool {
	return true
}

func (s str) Equals(object Object) bool {
	if object.Type() != object_type.Str {
		return false
	}
	return string(s) == object.Str()
}
