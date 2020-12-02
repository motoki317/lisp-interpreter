package node

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/token"
	"strings"
	"testing"
)

func (n *Node) equals(other *Node) bool {
	if n.Type != other.Type {
		fmt.Println("type mismatch")
		return false
	}
	switch n.Type {
	case Identifier:
		return n.Str == other.Str
	case Keyword:
		return n.Str == other.Str
	case Number:
		return n.Num == other.Num
	case Boolean:
		return n.B == other.B
	case Branch:
		if len(n.Children) != len(other.Children) {
			return false
		}
		for i := range n.Children {
			if !n.Children[i].equals(other.Children[i]) {
				return false
			}
		}
		return true
	}
	panic("unknown node type")
}

func allNodesEquals(first, second []*Node) bool {
	if len(first) != len(second) {
		return false
	}
	for i := range first {
		if !first[i].equals(second[i]) {
			return false
		}
	}
	return true
}

func readAllNodes(t *testing.T, p *Parser) []*Node {
	t.Helper()

	nodes := make([]*Node, 0)
	for {
		node, err := p.Next()
		if err == EOF {
			break
		}
		if err != nil {
			t.Fatalf("error while reading nodes: %v", err)
		}
		nodes = append(nodes, node)
	}
	return nodes
}

func TestParser(t *testing.T) {
	tests := []struct {
		name   string
		string string
		want   []*Node
	}{
		{
			name:   "single",
			string: "po po",
			want: []*Node{
				{Type: Identifier, Str: "po"},
				{Type: Identifier, Str: "po"},
			},
		},
		{
			name:   "simple node",
			string: "(define po -123 #t #f)",
			want: []*Node{
				{Type: Branch, Children: []*Node{
					{Type: Keyword, Str: "define"},
					{Type: Identifier, Str: "po"},
					{Type: Number, Num: -123},
					{Type: Boolean, B: true},
					{Type: Boolean, B: false},
				}},
			},
		},
		{
			name:   "nesting node",
			string: "(define (my-mult arg1 arg2) (* arg1 arg2))",
			want: []*Node{
				{Type: Branch, Children: []*Node{
					{Type: Keyword, Str: "define"},
					{Type: Branch, Children: []*Node{
						{Type: Identifier, Str: "my-mult"},
						{Type: Identifier, Str: "arg1"},
						{Type: Identifier, Str: "arg2"},
					}},
					{Type: Branch, Children: []*Node{
						{Type: Identifier, Str: "*"},
						{Type: Identifier, Str: "arg1"},
						{Type: Identifier, Str: "arg2"},
					}},
				}},
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			tokenizer := token.NewTokenizer(strings.NewReader(tt.string))
			parser := NewParser(tokenizer)
			if got := readAllNodes(t, parser); !allNodesEquals(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
