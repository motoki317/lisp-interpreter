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
}

func NewInterpreter(p *node.Parser, out io.Writer) *Interpreter {
	global := emptyFrame()
	for k, v := range defaultFuncs {
		global[k] = v
	}
	return &Interpreter{
		p:         p,
		out:       out,
		globalEnv: global,
	}
}

func (i *Interpreter) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(i.out, format, a...)
	if err != nil {
		fmt.Printf("Caught an error while writing to output: %v\n", err)
	}
}

func eval(n *node.Node, env env) *object {
	switch n.Type {
	case node.Keyword:
		return newErrorObject("unexpected keyword")
	case node.Identifier:
		if obj, ok := env.lookup(n.Str); ok {
			return obj
		} else {
			return newErrorObject(fmt.Sprintf("unbound identifier: %v", n.Str))
		}
	case node.Number:
		return newNumberObject(n.Num)
	}
	if n.Type != node.Branch {
		panic("node type not implemented")
	}

	// assert n.Type == node.Branch
	if len(n.Children) == 0 {
		return newErrorObject("bad syntax: empty sentence (\"()\") cannot be evaluated")
	}

	if n.Children[0].Type == node.Keyword {
		switch n.Children[0].Str {
		case "define":
			if len(n.Children) != 3 {
				return newErrorObject(fmt.Sprintf("bad syntax: define takes exactly 2 arguments, but got %v", len(n.Children)-1))
			}

			// define syntax sugar
			// (define (func-name arg1 arg2) ...)
			// = (define func-name (lambda (arg1 arg2) ...))
			if n.Children[1].Type == node.Branch {
				if len(n.Children[1].Children) == 0 || n.Children[1].Children[0].Type != node.Identifier {
					return newErrorObject("bad syntax: function definition requires function name")
				}

				// Rewrite AST and eval again
				funcName := n.Children[1].Children[0]
				argNames := n.Children[1].Children[1:]
				sentences := n.Children[2:]

				lambda := append([]*node.Node{
					{Type: node.Keyword, Str: "lambda"},
					{Type: node.Branch, Children: argNames},
				}, sentences...)
				return eval(&node.Node{Type: node.Branch, Children: []*node.Node{
					{Type: node.Keyword, Str: "define"},
					funcName,
					{Type: node.Branch, Children: lambda},
				},
				}, env)
			}

			// Normal define
			if n.Children[1].Type != node.Identifier {
				return newErrorObject(fmt.Sprintf("bad syntax: expected 1st argument of define to be identifier, but got %v", n.Children[1]))
			}

			key := n.Children[1].Str
			value := eval(n.Children[2], env)
			env.define(key, value)
			return voidObject
		case "lambda":
			if len(n.Children) < 3 {
				return newErrorObject(fmt.Sprintf("bad syntax: lambda takes 2 or more arguments, but got %v", len(n.Children)-1))
			}
			if n.Children[1].Type != node.Branch {
				return newErrorObject(fmt.Sprintf("bad syntax: 1st argument of lambda needs to be a list of arguments, but got %v", n.Children[1]))
			}

			argNames := make([]string, len(n.Children[1].Children))
			for i, arg := range n.Children[1].Children {
				if arg.Type != node.Identifier {
					return newErrorObject(fmt.Sprintf("bad syntax: expected %v-th argument of lambda function to be identifier, but got %v", i, arg.Type))
				}
				argNames[i] = arg.Str
			}
			sentences := n.Children[2:]
			return newFunctionObject(func(objects []*object) *object {
				if len(objects) != len(argNames) {
					return newErrorObject(fmt.Sprintf("expected length of arguments to be %v, but got %v", len(argNames), len(objects)))
				}

				newEnv := env.newEnv(newBindingFrame(argNames, objects))
				var ret *object
				for _, sentence := range sentences {
					ret = eval(sentence, newEnv)
				}
				return ret
			})
		}
	}

	objects := make([]*object, len(n.Children))
	for idx, child := range n.Children {
		objects[idx] = eval(child, env)
	}
	if objects[0].objectType != function {
		return newErrorObject(fmt.Sprintf("expected function in 0-th argument, but got %v", objects[0]))
	}
	return objects[0].f(objects[1:])
}

func (i *Interpreter) evalNext() (res *object, cont bool) {
	n, err := i.p.Next()
	if err == node.EOF {
		return nil, false
	}
	if err != nil {
		i.printf("An error occurred while parsing next input: %v\n", err)
	}
	return eval(n, newGlobalEnv(i.globalEnv)), true
}

func (i *Interpreter) ReadLoop() {
	for {
		i.printf("> ")
		res, cont := i.evalNext()
		if !cont {
			break
		}
		if res == voidObject {
			continue
		}
		i.printf("%v\n", res)
	}
}
