package node

import (
	"fmt"
	"strconv"
	"strings"
)

type Node struct {
	Type     Type
	Children []*Node
	Str      string
	Num      float64
	B        bool
}

type Type int

const (
	// Branch node has children nodes
	Branch Type = iota
	// Keyword Reserved keywords
	Keyword
	// Identifier Other strings
	Identifier
	// Number Numbers
	Number
	// Boolean Booleans
	Boolean
	// String String constant
	String
)

func (t Type) String() string {
	switch t {
	case Branch:
		return "branch"
	case Keyword:
		return "keyword"
	case Identifier:
		return "identifier"
	case Number:
		return "number"
	case String:
		return "string"
	}
	return strconv.Itoa(int(t))
}

func (n Node) String() string {
	switch n.Type {
	case Branch:
		formatted := make([]string, 0, len(n.Children))
		for _, child := range n.Children {
			formatted = append(formatted, child.String())
		}
		return "(" + strings.Join(formatted, " ") + ")"
	case Keyword:
		return n.Str
	case Identifier:
		return n.Str
	case Number:
		return fmt.Sprintf("%v", n.Num)
	case Boolean:
		if n.B {
			return "#t"
		} else {
			return "#f"
		}
	case String:
		return fmt.Sprintf("\"%v\"", n.Str)
	}
	return fmt.Sprintf("unknown_type: %v", n.Type)
}
