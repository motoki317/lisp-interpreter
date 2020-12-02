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
	// Word Other words
	Word
)

func (t Type) String() string {
	switch t {
	case LeftPar:
		return "left_par"
	case RightPar:
		return "right_par"
	case Word:
		return "word"
	}
	return strconv.Itoa(int(t))
}
