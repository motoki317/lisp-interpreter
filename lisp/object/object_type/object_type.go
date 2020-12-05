package object_type

import "strconv"

type T int

const (
	Number T = iota
	Boolean
	Symbol
	Str
	Cons
	Null
	Void
	Function
	Err
)

func (t T) String() string {
	switch t {
	case Number:
		return "number"
	case Boolean:
		return "boolean"
	case Symbol:
		return "symbol"
	case Str:
		return "string"
	case Cons:
		return "cons"
	case Null:
		return "null"
	case Void:
		return "void"
	case Function:
		return "function"
	case Err:
		return "error"
	}
	return strconv.Itoa(int(t))
}
