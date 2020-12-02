package lisp

import (
	"fmt"
	"strconv"
)

var (
	voidObject = &object{objectType: void}
	nullObject = &object{objectType: null}
)

type generalFunc func(objects []*object) *object

type object struct {
	objectType objectType
	num        float64
	b          bool
	pair       [2]*object
	str        string
	f          generalFunc
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

func newConsObject(car, cdr *object) *object {
	return &object{
		objectType: cons,
		pair:       [2]*object{car, cdr},
	}
}

func newFunctionObject(f generalFunc) *object {
	return &object{
		objectType: function,
		f:          f,
	}
}

func newErrorObject(msg string) *object {
	return &object{
		objectType: err,
		str:        msg,
	}
}
