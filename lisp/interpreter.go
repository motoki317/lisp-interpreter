package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
	"io"
)

type Interpreter struct {
	p         *node.Parser
	out       io.Writer
	globalEnv frame
	cuiMode   bool
}

func NewInterpreter(p *node.Parser, out io.Writer, cuiMode bool) *Interpreter {
	global := emptyFrame()
	for k, v := range defaultEnv {
		global[k] = v
	}
	i := &Interpreter{
		p:         p,
		out:       out,
		globalEnv: global,
		cuiMode:   cuiMode,
	}
	global["display"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			i.printf("%v\n", objects[0])
			return voidObject
		}))
	global["read"] = newFunctionObject(
		makeNullary(func(objects []*object) *object {
			n, err := i.p.Next()
			if err == node.EOF {
				return newErrorObject("end of input")
			}
			if err != nil {
				return newErrorObject(fmt.Sprintf("an error occurred while reading from input: %v", err))
			}
			return eval(&node.Node{
				Type: node.Branch,
				Children: []*node.Node{
					{Type: node.Keyword, Str: "quote"},
					n,
				},
			}, newGlobalEnv(i.globalEnv))
		}))
	return i
}

func (i *Interpreter) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(i.out, format, a...)
	if err != nil {
		fmt.Printf("Caught an error while writing to output: %v\n", err)
	}
}

func (i *Interpreter) evalNext() (res *object, cont bool) {
	n, err := i.p.Next()
	if err == node.EOF {
		return nil, false
	}
	if err != nil {
		i.printf("An error occurred while parsing next input: %v\n", err)
		return nil, true
	}
	return eval(n, newGlobalEnv(i.globalEnv)), true
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
		if res == nil || res == voidObject {
			continue
		}
		i.printf("%v\n", res)
	}
}
