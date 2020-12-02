package lisp

import (
	"bytes"
	"github.com/motoki317/lisp-interpreter/node"
	"github.com/motoki317/lisp-interpreter/token"
	"strings"
	"testing"
)

func TestInterpreter(t *testing.T) {
	tests := []struct {
		name    string
		inputs  []string
		outputs []string
	}{
		{
			name: "numbers",
			inputs: []string{
				"42",
				"334",
			},
			outputs: []string{
				"42\n",
				"334\n",
			},
		},
		{
			name: "basic arithmetic",
			inputs: []string{
				"(+ 1 2)",
				"(- 13 8)",
				"(* 15 20)",
				"(/ 300 50)",
				"(+ 1 2 (- 3 (* 4 5 (/ 10 5) 6) 7) 8 9)",
			},
			outputs: []string{
				"3\n",
				"5\n",
				"300\n",
				"6\n",
				"-224\n",
			},
		},
		{
			name: "define numbers",
			inputs: []string{
				"(define xx 2)",
				"(define po 5)",
				"(* xx po xx)",
			},
			outputs: []string{
				"",
				"",
				"20\n",
			},
		},
		{
			name: "basic lambda",
			inputs: []string{
				"(lambda (x) (* x 2))",
				"((lambda (x) (* x 2)) 2)",
			},
			outputs: []string{
				"<function>\n",
				"4\n",
			},
		},
		{
			name: "define lambda",
			inputs: []string{
				"(define double (lambda (x) (* x 2)))",
				"double",
				"(double 3)",
				"(double 5)",
			},
			outputs: []string{
				"",
				"<function>\n",
				"6\n",
				"10\n",
			},
		},
		{
			name: "define lambda (syntax sugar)",
			inputs: []string{
				"(define (double x) (* x 2))",
				"double",
				"(double 3)",
				"(double 5)",
			},
			outputs: []string{
				"",
				"<function>\n",
				"6\n",
				"10\n",
			},
		},
	}
	const waitOut = "> "
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			interpreter := NewInterpreter(node.NewParser(token.NewTokenizer(strings.NewReader(strings.Join(tt.inputs, "\n")))), out)
			interpreter.ReadLoop()

			expectOut := waitOut + strings.Join(tt.outputs, waitOut) + waitOut

			if gotOut := out.String(); gotOut != expectOut {
				t.Errorf("gotOut %v, want %v", gotOut, expectOut)
			}
		})
	}
}
