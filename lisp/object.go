package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
	"strconv"
)

var (
	voidObject = &object{objectType: void}
	nullObject = &object{objectType: null}
)

type tailOptFunc func(objects []*object) (*object, *node.Node, *env)
type generalFunc func(objects []*object) *object

type object struct {
	objectType objectType
	num        float64
	b          bool
	pair       [2]*object
	str        string
	f          tailOptFunc
}

func (o *object) equals(other *object) bool {
	if o.objectType != other.objectType {
		return false
	}
	switch o.objectType {
	case number:
		return o.num == other.num
	case boolean:
		return o.b == other.b
	case symbol:
		return o.str == other.str
	case str:
		return o.str == other.str
	case cons:
		return o.pair[0].equals(other.pair[0]) && o.pair[1].equals(other.pair[1])
	case null:
		return true
	case void:
		return true
	case function:
		return o == other
	case err:
		return o.str == other.str
	}
	panic("object type not implemented")
}

func (o *object) isList() bool {
	if o == nullObject {
		return true
	}
	return o.objectType == cons && o.pair[1].isList()
}

func (o *object) listElements() []*object {
	if o == nullObject {
		return []*object{}
	}
	first := o.pair[0]
	rest := o.pair[1].listElements()
	return append([]*object{first}, rest...)
}

func (o *object) stringStripPars() string {
	switch o.objectType {
	case cons:
		ret := o.pair[0].String()
		switch o.pair[1].objectType {
		case cons:
			ret += " " + o.pair[1].stringStripPars()
		case null:
		default:
			ret += " . " + o.pair[1].String()
		}
		return ret
	}
	return o.String()
}

func (o *object) String() string {
	switch o.objectType {
	case number:
		return fmt.Sprintf("%v", o.num)
	case boolean:
		if o.b {
			return "#t"
		} else {
			return "#f"
		}
	case symbol:
		return o.str
	case str:
		return fmt.Sprintf("\"%v\"", o.str)
	case cons:
		return "(" + o.stringStripPars() + ")"
	case null:
		return "()"
	case void:
		return "<void>"
	case function:
		return "<function>"
	case err:
		return fmt.Sprintf("error: %v", o.str)
	}
	panic("object type not implemented")
}

func (o *object) isTruthy() bool {
	return o.objectType != boolean || o.b
}

type objectType int

const (
	number objectType = iota
	boolean
	symbol
	str
	cons
	null
	void
	function
	err
)

func (t objectType) String() string {
	switch t {
	case number:
		return "number"
	case boolean:
		return "boolean"
	case symbol:
		return "symbol"
	case str:
		return "string"
	case cons:
		return "cons"
	case null:
		return "null"
	case void:
		return "void"
	case function:
		return "function"
	case err:
		return "error"
	}
	return strconv.Itoa(int(t))
}

func newNumberObject(num float64) *object {
	return &object{
		objectType: number,
		num:        num,
	}
}

func newBooleanObject(b bool) *object {
	return &object{
		objectType: boolean,
		b:          b,
	}
}

func newSymbolObject(str string) *object {
	return &object{
		objectType: symbol,
		str:        str,
	}
}

func newStringObject(s string) *object {
	return &object{
		objectType: str,
		str:        s,
	}
}

func newConsObject(car, cdr *object) *object {
	return &object{
		objectType: cons,
		pair:       [2]*object{car, cdr},
	}
}

func callWithTailOptimization(f func(objects []*object) (*object, *node.Node, *env), objects []*object) *object {
	obj, n, env := f(objects)
	if obj != nil {
		return obj
	}
	return evalWithTailOptimization(n, env)
}

func wrapFunction(f generalFunc) func(objects []*object) (*object, *node.Node, *env) {
	return func(objects []*object) (*object, *node.Node, *env) {
		return f(objects), nil, nil
	}
}

func newRawFunctionObject(f tailOptFunc) *object {
	return &object{
		objectType: function,
		f:          f,
	}
}

func newFunctionObject(f generalFunc) *object {
	return &object{
		objectType: function,
		f:          wrapFunction(f),
	}
}

func newErrorObject(msg string) *object {
	return &object{
		objectType: err,
		str:        msg,
	}
}
