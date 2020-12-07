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
		{
			name: "booleans",
			inputs: []string{
				"#t",
				"#f",
				"(> 3 2)",
				"(>= 3 3)",
				"(= 0 1)",
				"(zero? 0)",
				"(even? 4)",
				"(odd? 4)",
				"(not (= 0 1))",
			},
			outputs: []string{
				"#t\n",
				"#f\n",
				"#t\n",
				"#t\n",
				"#f\n",
				"#t\n",
				"#t\n",
				"#f\n",
				"#t\n",
			},
		},
		{
			name: "short circuit",
			inputs: []string{
				"(and)",
				"(or)",
				"(and 3 4)",
				"(and (= 5 0) (/ 5 0))",
				"(or #f #t)",
				"(or #f 5)",
			},
			outputs: []string{
				"#t\n",
				"#f\n",
				"4\n",
				"#f\n",
				"#t\n",
				"5\n",
			},
		},
		{
			name: "if",
			inputs: []string{
				"(define (my-div x y) (if (= y 0) 0 (/ x y)))",
				"(my-div 10 5)",
				"(my-div 10 0)",
			},
			outputs: []string{
				"",
				"2\n",
				"0\n",
			},
		},
		{
			name: "cond",
			inputs: []string{
				"(define (sign x)" +
					"(cond ((> x 0) 1)" +
					"      ((= x 0) 0)" +
					"      (else -1)))",
				"(sign 5)",
				"(sign 0)",
				"(sign -100)",
			},
			outputs: []string{
				"",
				"1\n",
				"0\n",
				"-1\n",
			},
		},
		{
			name: "let",
			inputs: []string{
				"(define (let-test x)" +
					"  (let ((x (+ x 1))" +
					"        (y (+ x 2)))" +
					"    (* x y)))",
				"(define (let-test-2 x)" +
					"  (let* ((x (+ x 1))" +
					"         (y (+ x 2)))" +
					"    (* x y)))",
				"(let-test 1)",
				"(let-test-2 1)",
			},
			outputs: []string{
				"",
				"",
				"6\n",
				"8\n",
			},
		},
		{
			name: "cons",
			inputs: []string{
				"(cons 1 2)",
				"(cons 1 (cons 2 3))",
				"(cons (cons 1 2) 3)",
				"(car (cons 1 2))",
				"(cdr (cons 1 2))",
				"(cadr (cons 1 (cons 2 3)))",
			},
			outputs: []string{
				"(1 . 2)\n",
				"(1 2 . 3)\n",
				"((1 . 2) . 3)\n",
				"1\n",
				"2\n",
				"2\n",
			},
		},
		{
			name: "quote",
			inputs: []string{
				"'po",
				"(quote po)",
				"'()",
				"'(1 2 3)",
				"(caddr '(1 2 3))",
				"'(1 . 2)",
				"(cdr '(1 . 2))",
				"'(define (xx po) (po))",
				"(cadadr '(define (xx po) (po)))",
			},
			outputs: []string{
				"po\n",
				"po\n",
				"()\n",
				"(1 2 3)\n",
				"3\n",
				"(1 . 2)\n",
				"2\n",
				"(define (xx po) (po))\n",
				"po\n",
			},
		},
		{
			name: "set!",
			inputs: []string{
				"(define po 20)",
				"po",
				"(set! po 50)",
				"po",
			},
			outputs: []string{
				"",
				"20\n",
				"",
				"50\n",
			},
		},
		{
			name: "set-car, cdr",
			inputs: []string{
				"(define p (cons 1 2))",
				"p",
				"(set-car! p 3)",
				"p",
				"(set-cdr! p 4)",
				"p",
			},
			outputs: []string{
				"",
				"(1 . 2)\n",
				"",
				"(3 . 2)\n",
				"",
				"(3 . 4)\n",
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			out := &bytes.Buffer{}
			interpreter := NewInterpreter(node.NewParser(token.NewTokenizer(strings.NewReader(strings.Join(tt.inputs, "\n")))), out, false, 0)
			interpreter.ReadLoop()

			expectOut := strings.Join(tt.outputs, "")

			if gotOut := out.String(); gotOut != expectOut {
				t.Errorf("gotOut %v, want %v", gotOut, expectOut)
			}
		})
	}
}
