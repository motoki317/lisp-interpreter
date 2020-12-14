package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

func NewConsObject(car, cdr Object) Object {
	return &cons{car, cdr}
}

func (c *cons) Type() object_type.T {
	return object_type.Cons
}

func (c *cons) Number() float64 {
	panic("number() called on cons object")
}

func (c *cons) Bool() bool {
	panic("Bool() called on cons object")
}

func (c *cons) Pair() *[2]Object {
	return (*[2]Object)(c)
}

func (c *cons) Str() string {
	panic("Str() called on cons object")
}

func (c *cons) F(_ []Object) (Object, *node.Node, *Env) {
	panic("F() called on cons object")
}

func (c *cons) stringStripPars(display bool) string {
	var ret string
	if display {
		ret = c[0].String()
	} else {
		ret = c[0].Display()
	}
	switch c[1].Type() {
	case object_type.Cons:
		ret += " " + (c[1]).(*cons).stringStripPars(display)
	case object_type.Null:
		// append none
	default:
		if display {
			ret += " . " + c[1].Display()
		} else {
			ret += " . " + c[1].String()
		}
	}
	return ret
}

func (c *cons) String() string {
	return "(" + c.stringStripPars(false) + ")"
}

func (c *cons) Display() string {
	return "(" + c.stringStripPars(true) + ")"
}

func (c *cons) IsList() bool {
	// this is cons type, only requires cdr to also be list
	return c[1].IsList()
}

func (c *cons) ListElements() []Object {
	if c.Type() == object_type.Null {
		return []Object{}
	}
	first := c[0]
	rest := c[1].ListElements()
	return append([]Object{first}, rest...)
}

func (c *cons) IsTruthy() bool {
	return true
}

func (c *cons) Equals(object Object) bool {
	if object.Type() != object_type.Cons {
		return false
	}
	o := object.(*cons)
	return c[0].Equals(o[0]) && c[1].Equals(o[1])
}
