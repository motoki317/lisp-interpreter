package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewPromiseObject(n *node.Node, e *Env) Object {
	return &promise{n: n, e: e}
}

func (d *promise) Type() object_type.T {
	return object_type.Promise
}

func (d *promise) Number() float64 {
	panic("Number() called on promise object")
}

func (d *promise) Bool() bool {
	panic("Bool() called on promise object")
}

func (d *promise) Pair() *[2]Object {
	panic("Pair() called on promise object")
}

func (d *promise) Str() string {
	panic("Str() called on promise object")
}

func (d *promise) F(_ []Object) (Object, *node.Node, *Env) {
	return nil, d.n, d.e
}

func (d *promise) String() string {
	return "<promise>"
}

func (d *promise) IsList() bool {
	return false
}

func (d *promise) ListElements() []Object {
	panic("ListElements() called on promise object")
}

func (d *promise) IsTruthy() bool {
	return true
}

func (d *promise) Equals(object Object) bool {
	if object.Type() != object_type.Promise {
		return false
	}
	o := object.(*promise)
	// not quite an accurate equality implementation
	return *d == *o
}
