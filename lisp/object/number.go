package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
	"strconv"
)

func NewNumberObject(num float64) Object {
	return number(num)
}

func (n number) Type() object_type.T {
	return object_type.Number
}

func (n number) Number() float64 {
	return float64(n)
}

func (n number) Bool() bool {
	panic("Bool() called on number object")
}

func (n number) Pair() *[2]Object {
	panic("Pair() called on number object")
}

func (n number) Str() string {
	panic("Str() called on number object")
}

func (n number) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on number object")
}

func (n number) String() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n number) Display() string {
	return strconv.FormatFloat(float64(n), 'f', -1, 64)
}

func (n number) IsList() bool {
	return false
}

func (n number) ListElements() []Object {
	panic("ListElements() called on number object")
}

func (n number) IsTruthy() bool {
	return true
}

func (n number) Equals(o Object) bool {
	if o.Type() != object_type.Number {
		return false
	}
	return float64(n) == o.Number()
}
