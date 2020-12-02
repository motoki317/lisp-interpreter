package node

import (
	"fmt"
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
)

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
	}
	return fmt.Sprintf("unknown_type: %v", n.Type)
}
