package macro_test

import (
	"github.com/motoki317/lisp-interpreter/lisp/macro"
	"github.com/motoki317/lisp-interpreter/lisp/object"
	"github.com/motoki317/lisp-interpreter/node"
	"github.com/motoki317/lisp-interpreter/token"
	"reflect"
	"strings"
	"testing"
)

func read(t *testing.T, input string) *node.Node {
	t.Helper()

	p := node.NewParser(token.NewTokenizer(strings.NewReader(input)))
	n, err := p.Next()
	if err != nil {
		t.Fatalf("error when reading input: %v", err)
	}
	return n
}

func TestMacro_Replace(t *testing.T) {
	// http://www.shido.info/lisp/scheme_syntax_e.html
	tests := []struct {
		name   string
		macros []string
		input  string
		want   *node.Node
	}{
		{
			name: "basic",
			macros: []string{
				"(define-syntax nil! (syntax-rules () ((_ x) (set! x '()))))",
			},
			input: "(nil! x)",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "set!"},
					{Type: node.Identifier, Str: "x"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "quote"},
						{Type: node.Branch, Children: []*node.Node{}},
					}},
				},
			},
		},
		{
			name: "variadic capture 1",
			macros: []string{
				"(define-syntax when (syntax-rules () ((_ pred b1 ...) (if pred (begin b1 ...)))))",
			},
			input: "(when #t b1 b2 b3)",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "if"},
					{Type: node.Boolean, B: true},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "begin"},
						{Type: node.Identifier, Str: "b1"},
						{Type: node.Identifier, Str: "b2"},
						{Type: node.Identifier, Str: "b3"},
					}},
				},
			},
		},
		{
			name: "variadic capture 2",
			macros: []string{
				"(define-syntax when (syntax-rules () ((_ pred b1 ...) (if pred (begin b1 ...)))))",
			},
			input: "(when my-pred b1)",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "if"},
					{Type: node.Identifier, Str: "my-pred"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "begin"},
						{Type: node.Identifier, Str: "b1"},
					}},
				},
			},
		},
		{
			name: "internal application",
			macros: []string{
				"(define-syntax when (syntax-rules () ((_ pred b1 ...) (if pred (begin b1 ...)))))",
			},
			input: "(define (f x) (when (= x 0) (display \"zero\")) (- x 1))",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "define"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Identifier, Str: "f"},
						{Type: node.Identifier, Str: "x"},
					}},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "if"},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Identifier, Str: "="},
							{Type: node.Identifier, Str: "x"},
							{Type: node.Number, Num: 0},
						}},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Keyword, Str: "begin"},
							{Type: node.Branch, Children: []*node.Node{
								{Type: node.Identifier, Str: "display"},
								{Type: node.String, Str: "zero"},
							}},
						}},
					}},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Identifier, Str: "-"},
						{Type: node.Identifier, Str: "x"},
						{Type: node.Number, Num: 1},
					}},
				},
			},
		},
		{
			name: "recursive application",
			macros: []string{
				"(define-syntax when (syntax-rules () ((_ pred b1 ...) (if pred (begin b1 ...)))))",
				"(define-syntax while (syntax-rules () ((_ pred b1 ...) (begin (define (loop) (when pred b1 ... (loop))) (loop)))))",
			},
			input: "(while (< i 10) (display i) (display \" \") (set! x (+ x 1)))",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "begin"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "define"},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Identifier, Str: "loop"},
						}},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Keyword, Str: "if"},
							{Type: node.Branch, Children: []*node.Node{
								{Type: node.Identifier, Str: "<"},
								{Type: node.Identifier, Str: "i"},
								{Type: node.Number, Num: 10},
							}},
							{Type: node.Branch, Children: []*node.Node{
								{Type: node.Keyword, Str: "begin"},
								{Type: node.Branch, Children: []*node.Node{
									{Type: node.Identifier, Str: "display"},
									{Type: node.Identifier, Str: "i"},
								}},
								{Type: node.Branch, Children: []*node.Node{
									{Type: node.Identifier, Str: "display"},
									{Type: node.String, Str: " "},
								}},
								{Type: node.Branch, Children: []*node.Node{
									{Type: node.Keyword, Str: "set!"},
									{Type: node.Identifier, Str: "x"},
									{Type: node.Branch, Children: []*node.Node{
										{Type: node.Identifier, Str: "+"},
										{Type: node.Identifier, Str: "x"},
										{Type: node.Number, Num: 1},
									}},
								}},
								{Type: node.Branch, Children: []*node.Node{
									{Type: node.Identifier, Str: "loop"},
								}},
							}},
						}},
					}},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Identifier, Str: "loop"},
					}},
				},
			},
		},
		{
			name: "multiple branches 1",
			macros: []string{
				"(define-syntax inc! (syntax-rules () ((_ x) (begin (set! x (+ x 1)) x)) ((_ x i) (begin (set! x (+ x i)) x))))",
			},
			input: "(inc! i)",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "begin"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "set!"},
						{Type: node.Identifier, Str: "i"},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Identifier, Str: "+"},
							{Type: node.Identifier, Str: "i"},
							{Type: node.Number, Num: 1},
						}},
					}},
					{Type: node.Identifier, Str: "i"},
				},
			},
		},
		{
			name: "multiple branches 1",
			macros: []string{
				"(define-syntax inc! (syntax-rules () ((_ x) (begin (set! x (+ x 1)) x)) ((_ x i) (begin (set! x (+ x i)) x))))",
			},
			input: "(inc! i 3)",
			want: &node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "begin"},
					{Type: node.Branch, Children: []*node.Node{
						{Type: node.Keyword, Str: "set!"},
						{Type: node.Identifier, Str: "i"},
						{Type: node.Branch, Children: []*node.Node{
							{Type: node.Identifier, Str: "+"},
							{Type: node.Identifier, Str: "i"},
							{Type: node.Number, Num: 3},
						}},
					}},
					{Type: node.Identifier, Str: "i"},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			e := object.NewGlobalEnv(object.EmptyFrame())

			for _, macroStr := range tt.macros {
				macroCode := read(t, macroStr)
				m, err := macro.NewMacro(macroCode)
				if err != nil {
					t.Fatalf("error when creating macro: %v", err)
				}
				e.DefineGlobalMacro(m)
			}

			inputCode := read(t, tt.input)
			if got, err := e.ApplyMacro(inputCode); err != nil {
				t.Fatalf("error when applying macro: %v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v, want %v", got, tt.want)
			}
		})
	}
}
