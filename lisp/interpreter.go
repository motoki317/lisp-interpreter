package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/lisp/object"
	"github.com/motoki317/lisp-interpreter/node"
	"github.com/motoki317/lisp-interpreter/token"
	"io"
	"time"
)

type Interpreter struct {
	p         *node.Parser
	out       io.Writer
	globalEnv object.Frame
	cuiMode   bool
	timeout   time.Duration
}

func NewInterpreter(p *node.Parser, out io.Writer, cuiMode bool, timeout time.Duration) *Interpreter {
	global := object.EmptyFrame()
	for k, v := range defaultEnv {
		global[k] = v
	}
	i := &Interpreter{
		p:         p,
		out:       out,
		globalEnv: global,
		cuiMode:   cuiMode,
		timeout:   timeout,
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

// SetTokenizer sets internal tokenizer used by parser, to start using from the next call.
func (i *Interpreter) SetTokenizer(t *token.Tokenizer) {
	i.p.SetTokenizer(t)
}

// SetOutput sets the output used by this interpreter.
func (i *Interpreter) SetOutput(out io.Writer) {
	i.out = out
}

func (i *Interpreter) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(i.out, format, a...)
	if err != nil {
		fmt.Printf("Caught an error while writing to output: %v\n", err)
	}
}

func (i *Interpreter) evalNext() (res object.Object, cont bool, timedOut bool) {
	n, err := i.p.Next()
	if err == node.EOF {
		return nil, false, false
	}
	if err != nil {
		i.printf("An error occurred while parsing next input: %v\n", err)
		return nil, true, false
	}

	var stopper <-chan time.Time
	if i.timeout != time.Duration(0) {
		timer := time.NewTimer(i.timeout)
		stopper = timer.C
		defer timer.Stop()
	}

	res, timedOut = evalWithStopper(n, object.NewGlobalEnv(i.globalEnv), stopper)
	if timedOut {
		return nil, true, true
	} else {
		return res, true, false
	}
}

// ReadLoop executes the Read, Eval, Print loop (REPL), until the parser hits EOF.
func (i *Interpreter) ReadLoop() {
	for {
		if i.cuiMode {
			i.printf("> ")
		}
		res, cont, timedOut := i.evalNext()
		if !cont {
			break
		}
		if timedOut {
			i.printf("Timed out.\n")
			continue
		}
		if res == nil || res == object.VoidObj {
			continue
		}
		i.printf("%v\n", res)
	}
}
