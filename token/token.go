package token

import (
	"strconv"
)

type Token struct {
	Type   Type
	String string
}

type Type int

const (
	// LeftPar (
	LeftPar Type = iota
	// RightPar )
	RightPar
	// Keyword Reserved keywords: define
	Keyword
	// Identifier Other strings
	Identifier
	// Number Numbers
	Number
)

func (t Type) String() string {
	switch t {
	case LeftPar:
		return "left_par"
	case RightPar:
		return "right_par"
	case Keyword:
		return "keyword"
	case Identifier:
		return "identifier"
	case Number:
		return "number"
	}
	return strconv.Itoa(int(t))
}
