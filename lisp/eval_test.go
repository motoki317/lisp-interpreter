package lisp

import (
	"bytes"
	"fmt"
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
	"github.com/motoki317/lisp-interpreter/token"
	"strconv"
	"testing"
)

func setUpInterpreter(b *testing.B, preEval string) (*bytes.Buffer, *Interpreter) {
	b.Helper()

	input := bytes.NewBufferString("")
	parser := node.NewParser(token.NewTokenizer(input))
	out := &bytes.Buffer{}
	interpreter := NewInterpreter(parser, out, false, 0)

	input.WriteString(preEval)
	_, cont, _ := interpreter.evalNext()
	if !cont {
		b.Fatalf("not continued")
	}

	return input, interpreter
}

func BenchmarkEvalSumN(b *testing.B) {
	input, interpreter := setUpInterpreter(b, "(define (sum n) (if (<= n 0) 0 (+ n (sum (- n 1)))))")

	input.WriteString("(sum " + strconv.Itoa(b.N) + ")")
	b.ResetTimer()
	obj, cont, _ := interpreter.evalNext()
	b.StopTimer()
	if !cont {
		panic("not continued")
	}
	if obj.Type() != object_type.Number || int(obj.Number()) != b.N*(b.N+1)/2 {
		panic(fmt.Sprintf("unexpected object: %v", obj))
	}
}

func BenchmarkEvalSumTailN(b *testing.B) {
	input, interpreter := setUpInterpreter(b, "(define (sum-tail n a) (if (<= n 0) a (sum-tail (- n 1) (+ n a))))")

	input.WriteString("(sum-tail " + strconv.Itoa(b.N) + " 0)")
	b.ResetTimer()
	obj, cont, _ := interpreter.evalNext()
	b.StopTimer()
	if !cont {
		panic("not continued")
	}
	if obj.Type() != object_type.Number || int(obj.Number()) != b.N*(b.N+1)/2 {
		panic(fmt.Sprintf("unexpected object: %v", obj))
	}
}

func BenchmarkEvalSum(b *testing.B) {
	input, interpreter := setUpInterpreter(b, "(define (sum n) (if (<= n 0) 0 (+ n (sum (- n 1)))))")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input.WriteString("(sum 10000)")
		b.StartTimer()
		obj, cont, _ := interpreter.evalNext()
		b.StopTimer()
		if !cont {
			panic("not continued")
		}
		if obj.Type() != object_type.Number || obj.Number() != 50005000 {
			panic(fmt.Sprintf("unexpected object: %v", obj))
		}
	}
}

func BenchmarkEvalSumTail(b *testing.B) {
	input, interpreter := setUpInterpreter(b, "(define (sum-tail n a) (if (<= n 0) a (sum-tail (- n 1) (+ n a))))")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input.WriteString("(sum-tail 10000 0)")
		b.StartTimer()
		obj, cont, _ := interpreter.evalNext()
		b.StopTimer()
		if !cont {
			panic("not continued")
		}
		if obj.Type() != object_type.Number || obj.Number() != 50005000 {
			panic(fmt.Sprintf("unexpected object: %v", obj))
		}
	}
}
