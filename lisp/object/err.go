package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewErrorObject(msg string) Object {
	return err(msg)
}

func (e err) Type() object_type.T {
	return object_type.Err
}

func (e err) Number() float64 {
	panic("number() called on err object")
}

func (e err) Bool() bool {
	panic("Bool() called on err object")
}

func (e err) Pair() *[2]Object {
	panic("Pair() called on err object")
}

func (e err) Str() string {
	return string(e)
}

func (e err) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on err object")
}

func (e err) String() string {
	return "error: " + string(e)
}

func (e err) Display() string {
	return "error: " + string(e)
}

func (e err) IsList() bool {
	return false
}

func (e err) ListElements() []Object {
	panic("ListElements() called on err object")
}

func (e err) IsTruthy() bool {
	return true
}

func (e err) Equals(object Object) bool {
	if object.Type() != object_type.Err {
		return false
	}
	return string(e) == object.Str()
}
