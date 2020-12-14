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
				"42",
				"334",
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
				"3",
				"5",
				"300",
				"6",
				"-224",
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
				"20",
			},
		},
		{
			name: "basic lambda",
			inputs: []string{
				"(lambda (x) (* x 2))",
				"((lambda (x) (* x 2)) 2)",
			},
			outputs: []string{
				"<function>",
				"4",
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
				"<function>",
				"6",
				"10",
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
				"<function>",
				"6",
				"10",
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
				"#t",
				"#f",
				"#t",
				"#t",
				"#f",
				"#t",
				"#t",
				"#f",
				"#t",
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
				"#t",
				"#f",
				"4",
				"#f",
				"#t",
				"5",
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
				"2",
				"0",
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
				"1",
				"0",
				"-1",
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
				"6",
				"8",
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
				"(1 . 2)",
				"(1 2 . 3)",
				"((1 . 2) . 3)",
				"1",
				"2",
				"2",
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
				"po",
				"po",
				"()",
				"(1 2 3)",
				"3",
				"(1 . 2)",
				"2",
				"(define (xx po) (po))",
				"po",
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
				"20",
				"50",
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
				"(1 . 2)",
				"(3 . 2)",
				"(3 . 4)",
			},
		},
		{
			name: "variadic length functions 1",
			inputs: []string{
				"(define f (lambda x x))",
				"(f)",
				"(f 1)",
				"(f 1 2 3 4 5)",
				"(define (f . x) x)",
				"(f)",
				"(f 1)",
				"(f 1 2 3 4 5)",
			},
			outputs: []string{
				"()",
				"(1)",
				"(1 2 3 4 5)",
				"()",
				"(1)",
				"(1 2 3 4 5)",
			},
		},
		{
			name: "variadic length functions 2",
			inputs: []string{
				"(define f (lambda (x y . z) (list x y z)))",
				"(f 1 2)",
				"(f 1 2 3)",
				"(f 1 2 3 4 5)",
				"(define (f x y . z) (list x y z))",
				"(f 1 2)",
				"(f 1 2 3)",
				"(f 1 2 3 4 5)",
			},
			outputs: []string{
				"(1 2 ())",
				"(1 2 (3))",
				"(1 2 (3 4 5))",
				"(1 2 ())",
				"(1 2 (3))",
				"(1 2 (3 4 5))",
			},
		},
		{
			// http://www.shido.info/lisp/scheme_syntax_e.html
			name: "macros",
			inputs: []string{
				"(define-syntax when (syntax-rules () ((_ pred b1 ...) (if pred (begin b1 ...)))))",
				"(define-syntax while (syntax-rules () ((_ pred b1 ...) (begin (define (loop) (when pred b1 ... (loop))) (loop)))))",
				"(define-syntax for (syntax-rules () ((_ (i from to) b1 ...) (begin (define (loop i) (when (< i to) b1 ... (loop (+ i 1)))) (loop from)))))",
				"(define-syntax inc! (syntax-rules () ((_ x) (begin (set! x (+ x 1)) x)) ((_ x i) (begin (set! x (+ x i)) x))))",
				"(when #f (/ 1 0))", // no output
				"(let ((i 0)) (while (< i 3) (display i) (set! i (+ i 1))))", // 0, 1, 2
				"(for (i 0 3) (display i))",                                  // 0, 1, 2
				"(define i 0)",
				"(inc! i)",   // 1
				"(inc! i 3)", // 4
				"i",          // 4
			},
			outputs: []string{
				"0",
				"1",
				"2",
				"0",
				"1",
				"2",
				"1",
				"4",
				"4",
			},
		},
		{
			name: "promise / stream",
			inputs: []string{
				"(define-syntax s-cons (syntax-rules () ((_ a b) (cons a (delay b)))))",
				"(define (s-car s) (car s))",
				"(define (s-cdr s) (force (cdr s)))",
				"(define (s-null? s) (null? s))",
				"(define (s-head s n) (cond ((s-null? s) '()) ((<= n 0) '()) (else (cons (s-car s) (s-head (s-cdr s) (- n 1))))))",
				"(define (integers-from n) (s-cons n (integers-from (+ n 1))))",
				"(define integers* (integers-from 1))",
				"(s-head integers* 10)",
				"(s-head integers* 5)",
			},
			outputs: []string{
				"(1 2 3 4 5 6 7 8 9 10)",
				"(1 2 3 4 5)",
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

			expectOut := strings.Join(tt.outputs, "\n") + "\n"

			if gotOut := out.String(); gotOut != expectOut {
				t.Errorf("gotOut %v, want %v", gotOut, expectOut)
			}
		})
	}
}
