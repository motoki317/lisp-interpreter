package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewBooleanObject(b bool) Object {
	return boolean(b)
}

func (b boolean) Type() object_type.T {
	return object_type.Boolean
}

func (b boolean) Number() float64 {
	panic("number() called on boolean object")
}

func (b boolean) Bool() bool {
	return bool(b)
}

func (b boolean) Pair() *[2]Object {
	panic("Pair() called on boolean object")
}

func (b boolean) Str() string {
	panic("Str() called on boolean object")
}

func (b boolean) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on boolean object")
}

func (b boolean) String() string {
	if b {
		return "#t"
	} else {
		return "#f"
	}
}

func (b boolean) Display() string {
	if b {
		return "#t"
	} else {
		return "#f"
	}
}

func (b boolean) IsList() bool {
	return false
}

func (b boolean) ListElements() []Object {
	panic("ListElements() called on boolean object")
}

func (b boolean) IsTruthy() bool {
	return bool(b)
}

func (b boolean) Equals(object Object) bool {
	if object.Type() != object_type.Boolean {
		return false
	}
	return bool(b) == object.Bool()
}
