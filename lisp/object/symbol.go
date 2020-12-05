package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewSymbolObject(str string) Object {
	return symbol(str)
}

func (s symbol) Type() object_type.T {
	return object_type.Symbol
}

func (s symbol) Number() float64 {
	panic("number() called on symbol object")
}

func (s symbol) Bool() bool {
	panic("Bool() called on symbol object")
}

func (s symbol) Pair() *[2]Object {
	panic("Pair() called on symbol object")
}

func (s symbol) Str() string {
	return string(s)
}

func (s symbol) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on symbol object")
}

func (s symbol) String() string {
	return string(s)
}

func (s symbol) IsList() bool {
	return false
}

func (s symbol) ListElements() []Object {
	panic("ListElements() called on symbol object")
}

func (s symbol) IsTruthy() bool {
	return true
}

func (s symbol) Equals(object Object) bool {
	if object.Type() != object_type.Symbol {
		return false
	}
	return string(s) == object.Str()
}
