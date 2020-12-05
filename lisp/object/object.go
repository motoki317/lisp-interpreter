package object

import (
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
)

var (
	VoidObj = void{}
	NullObj = null{}
)

type Object interface {
	Type() object_type.T

	Number() float64
	Bool() bool
	Pair() *[2]Object
	// Str returns string data if type is symbol of Str.
	// Panics otherwise.
	Str() string
	F(objects []Object) (Object, *node.Node, *Env)

	String() string

	IsList() bool
	ListElements() []Object
	IsTruthy() bool
	Equals(Object) bool
}

type (
	number   float64
	boolean  bool
	symbol   string
	str      string
	cons     [2]Object
	null     struct{}
	void     struct{}
	function func(objects []Object) (Object, *node.Node, *Env)
	err      string
)
