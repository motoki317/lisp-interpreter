package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
)

func evalAnd(n *node.Node, env env) *object {
	// Short circuit evaluation
	res := newBooleanObject(true)
	for _, child := range n.Children[1:] {
		res = eval(child, env)
		if !res.isTruthy() {
			return newBooleanObject(false)
		}
	}
	return res
}

func evalOr(n *node.Node, env env) *object {
	// Short circuit evaluation
	for _, child := range n.Children[1:] {
		res := eval(child, env)
		if res.isTruthy() {
			return res
		}
	}
	return newBooleanObject(false)
}

func evalIf(n *node.Node, env env) *object {
	if len(n.Children) != 3 && len(n.Children) != 4 {
		return newErrorObject(fmt.Sprintf("bad syntax: if needs 2 or 3 arguments, but got %v", len(n.Children)-1))
	}

	res := eval(n.Children[1], env)
	if res.isTruthy() {
		return eval(n.Children[2], env)
	} else {
		if len(n.Children) == 4 {
			return eval(n.Children[3], env)
		} else {
			return voidObject
		}
	}
}

func evalLet(n *node.Node, env env) *object {
	if len(n.Children) <= 2 {
		return newErrorObject(fmt.Sprintf("bad syntax: let needs at least 2 arguments, but got %v", len(n.Children)-1))
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	keys := make([]string, len(pairs.Children))
	values := make([]*object, len(pairs.Children))
	for i, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return newErrorObject(fmt.Sprintf("bad syntax: let bind pair needs a list of length 2, but got length %v", len(pair.Children)))
		}
		if pair.Children[0].Type != node.Identifier {
			return newErrorObject(fmt.Sprintf("bad syntax: let bind pair requires identifier, but got %v", pair.Children[0].Type))
		}

		keys[i] = pair.Children[0].Str
		values[i] = eval(pair.Children[1], env)
	}

	env = env.newEnv(newBindingFrame(keys, values))
	res := voidObject
	for _, sentence := range sentences {
		res = eval(sentence, env)
	}
	return res
}

func evalLetSeq(n *node.Node, env env) *object {
	if len(n.Children) <= 2 {
		return newErrorObject(fmt.Sprintf("bad syntax: let* needs at least 2 arguments, but got %v", len(n.Children)-1))
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	env = env.newEnv(emptyFrame())
	for _, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return newErrorObject(fmt.Sprintf("bad syntax: let* bind pair needs a list of length 2, but got length %v", len(pair.Children)))
		}
		if pair.Children[0].Type != node.Identifier {
			return newErrorObject(fmt.Sprintf("bad syntax: let* bind pair requires identifier, but got %v", pair.Children[0].Type))
		}

		key := pair.Children[0].Str
		value := eval(pair.Children[1], env)

		env.define(key, value)
	}

	res := voidObject
	for _, sentence := range sentences {
		res = eval(sentence, env)
	}
	return res
}

func evalCond(n *node.Node, env env) *object {
	if len(n.Children) == 1 {
		return newErrorObject("bad syntax: cond needs at least 1 argument, but got 0")
	}

	for _, branch := range n.Children[1:] {
		if branch.Type != node.Branch || len(branch.Children) == 0 {
			return newErrorObject("bad syntax: cond bad branch")
		}

		test := branch.Children[0]
		if (test.Type == node.Keyword && test.Str == "else") || eval(test, env).isTruthy() {
			res := voidObject
			for _, child := range branch.Children[1:] {
				res = eval(child, env)
			}
			return res
		}
	}
	return voidObject
}

func evalQuote(n *node.Node) *object {
	switch n.Type {
	case node.Number:
		return newNumberObject(n.Num)
	case node.Boolean:
		return newBooleanObject(n.B)
	case node.Identifier:
		return newSymbolObject(n.Str)
	case node.Keyword:
		return newSymbolObject(n.Str)
	}
	if n.Type != node.Branch {
		panic(fmt.Sprintf("quote node type not implemented: %v", n.Type))
	}
	// assert n.Type == node.Branch
	if len(n.Children) == 0 {
		return nullObject
	}
	if len(n.Children) == 3 && n.Children[1].Type == node.Identifier && n.Children[1].Str == "." {
		return newConsObject(
			evalQuote(n.Children[0]),
			evalQuote(n.Children[2]))
	}
	return newConsObject(
		evalQuote(n.Children[0]),
		evalQuote(&node.Node{
			Type:     node.Branch,
			Children: n.Children[1:],
		}))
}

func evalDefine(n *node.Node, env env) *object {
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
	if len(n.Children) != 3 {
		return newErrorObject(fmt.Sprintf("bad syntax: define takes exactly 2 arguments, but got %v", len(n.Children)-1))
	}

	if n.Children[1].Type != node.Identifier {
		return newErrorObject(fmt.Sprintf("bad syntax: expected 1st argument of define to be identifier, but got %v", n.Children[1]))
	}

	key := n.Children[1].Str
	value := eval(n.Children[2], env)
	env.define(key, value)
	return voidObject
}

func evalLambda(n *node.Node, env env) *object {
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
	case node.Boolean:
		return newBooleanObject(n.B)
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
		case "and":
			return evalAnd(n, env)
		case "or":
			return evalOr(n, env)
		case "if":
			return evalIf(n, env)
		case "let":
			return evalLet(n, env)
		case "let*":
			return evalLetSeq(n, env)
		case "cond":
			return evalCond(n, env)
		case "quote":
			if len(n.Children) != 2 {
				return newErrorObject(fmt.Sprintf("quote needs exactly 1 argument, but got %v", len(n.Children)-1))
			}
			return evalQuote(n.Children[1])
		case "define":
			return evalDefine(n, env)
		case "lambda":
			return evalLambda(n, env)
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
