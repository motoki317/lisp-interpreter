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
	}
	return fmt.Sprintf("unknown_type: %v", n.Type)
}
