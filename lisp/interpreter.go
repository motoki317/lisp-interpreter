package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/lisp/object"
	"github.com/motoki317/lisp-interpreter/node"
	"io"
)

type Interpreter struct {
	p         *node.Parser
	out       io.Writer
	globalEnv object.Frame
	cuiMode   bool
}

func NewInterpreter(p *node.Parser, out io.Writer, cuiMode bool) *Interpreter {
	global := object.EmptyFrame()
	for k, v := range defaultEnv {
		global[k] = v
	}
	i := &Interpreter{
		p:         p,
		out:       out,
		globalEnv: global,
		cuiMode:   cuiMode,
	}
	global["display"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			i.printf("%v\n", objects[0])
			return object.VoidObj
		}))
	global["read"] = object.NewWrappedFunctionObject(
		makeNullary(func(objects []object.Object) object.Object {
			n, err := i.p.Next()
			if err == node.EOF {
				return object.NewErrorObject("end of input")
			}
			if err != nil {
				return object.NewErrorObject(fmt.Sprintf("an error occurred while reading from input: %v", err))
			}
			return evalWithTailOptimization(&node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "quote"},
					n,
				},
			}, object.NewGlobalEnv(i.globalEnv))
		}))
	return i
}

func (i *Interpreter) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(i.out, format, a...)
	if err != nil {
		fmt.Printf("Caught an error while writing to output: %v\n", err)
	}
}

func (i *Interpreter) evalNext() (res object.Object, cont bool) {
	n, err := i.p.Next()
	if err == node.EOF {
		return nil, false
	}
	if err != nil {
		i.printf("An error occurred while parsing next input: %v\n", err)
		return nil, true
	}
	return evalWithTailOptimization(n, object.NewGlobalEnv(i.globalEnv)), true
}

func (i *Interpreter) ReadLoop() {
	for {
		if i.cuiMode {
			i.printf("> ")
		}
		res, cont := i.evalNext()
		if !cont {
			break
		}
		if res == nil || res == object.VoidObj {
			continue
		}
		i.printf("%v\n", res)
	}
}
