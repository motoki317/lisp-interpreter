package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
	"io"
)

type Interpreter struct {
	p         *node.Parser
	out       io.Writer
	globalEnv map[string]*object
}

func NewInterpreter(p *node.Parser, out io.Writer) *Interpreter {
	globalEnv := make(map[string]*object)
	for str, obj := range defaultFuncs {
		globalEnv[str] = obj
	}
	return &Interpreter{
		p:         p,
		out:       out,
		globalEnv: globalEnv,
	}
}

func (i *Interpreter) printf(format string, a ...interface{}) {
	_, err := fmt.Fprintf(i.out, format, a...)
	if err != nil {
		fmt.Printf("Caught an error while writing to output: %v\n", err)
	}
}

func (i *Interpreter) eval(n *node.Node) *object {
	switch n.Type {
	case node.Keyword:
		return newErrorObject("unexpected keyword")
	case node.Identifier:
		if obj, ok := i.globalEnv[n.Str]; ok {
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
		return newErrorObject("empty sentence (\"()\") cannot be evaluated")
	}

	if n.Children[0].Type == node.Keyword {
		switch n.Children[0].Str {
		case "define":
			if n.Children[1].Type != node.Identifier {
				return newErrorObject(fmt.Sprintf("expected 1st argument of define to be identifier, but got %v", n.Children[1]))
			}
			name := n.Children[1].Str
			val := i.eval(n.Children[2])
			i.globalEnv[name] = val
			return voidObject
		}
	}

	objects := make([]*object, len(n.Children))
	for idx, child := range n.Children {
		objects[idx] = i.eval(child)
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
	return i.eval(n), true
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
