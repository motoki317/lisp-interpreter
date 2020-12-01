package lisp

import (
	"fmt"
	"strconv"
)

var (
	voidObject = &object{objectType: void}
)

type generalFunc func(objects []*object) *object

type object struct {
	objectType objectType
	num        float64
	str        string
	f          generalFunc
}

func (o *object) String() string {
	switch o.objectType {
	case number:
		return fmt.Sprintf("%v", o.num)
	case void:
		return "<void>"
	case function:
		return "<function>"
	case err:
		return fmt.Sprintf("error: %v", o.str)
	}
	panic("object type not implemented")
}

type objectType int

const (
	number objectType = iota
	void
	function
	err
)

func (t objectType) String() string {
	switch t {
	case number:
		return "number"
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

func newErrorObject(msg string) *object {
	return &object{
		objectType: err,
		str:        msg,
	}
}
