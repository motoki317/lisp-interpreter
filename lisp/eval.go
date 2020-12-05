package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
)

func evalAnd(n *node.Node, env *env) *object {
	// Short circuit evaluation
	res := newBooleanObject(true)
	for _, child := range n.Children[1:] {
		res = evalWithTailOptimization(child, env)
		if !res.isTruthy() {
			return newBooleanObject(false)
		}
	}
	return res
}

func evalOr(n *node.Node, env *env) *object {
	// Short circuit evaluation
	for _, child := range n.Children[1:] {
		res := evalWithTailOptimization(child, env)
		if res.isTruthy() {
			return res
		}
	}
	return newBooleanObject(false)
}

func evalIf(n *node.Node, env *env) (*object, *node.Node, *env) {
	if len(n.Children) != 3 && len(n.Children) != 4 {
		return newErrorObject(fmt.Sprintf("bad syntax: if needs 2 or 3 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	res := evalWithTailOptimization(n.Children[1], env)
	if res.isTruthy() {
		return nil, n.Children[2], env
	} else {
		if len(n.Children) == 4 {
			return nil, n.Children[3], env
		} else {
			return voidObject, nil, nil
		}
	}
}

func evalLet(n *node.Node, env *env) (*object, *node.Node, *env) {
	if len(n.Children) <= 2 {
		return newErrorObject(fmt.Sprintf("bad syntax: let needs at least 2 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	keys := make([]string, len(pairs.Children))
	values := make([]*object, len(pairs.Children))
	for i, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return newErrorObject(fmt.Sprintf("bad syntax: let bind pair needs a list of length 2, but got length %v", len(pair.Children))), nil, nil
		}
		if pair.Children[0].Type != node.Identifier {
			return newErrorObject(fmt.Sprintf("bad syntax: let bind pair requires identifier, but got %v", pair.Children[0].Type)), nil, nil
		}

		keys[i] = pair.Children[0].Str
		values[i] = evalWithTailOptimization(pair.Children[1], env)
	}

	env = env.newEnv(newBindingFrame(keys, values))
	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, env)
	}
	return nil, sentences[len(sentences)-1], env
}

func evalLetSeq(n *node.Node, env *env) (*object, *node.Node, *env) {
	if len(n.Children) <= 2 {
		return newErrorObject(fmt.Sprintf("bad syntax: let* needs at least 2 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	env = env.newEnv(emptyFrame())
	for _, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return newErrorObject(fmt.Sprintf("bad syntax: let* bind pair needs a list of length 2, but got length %v", len(pair.Children))), nil, nil
		}
		if pair.Children[0].Type != node.Identifier {
			return newErrorObject(fmt.Sprintf("bad syntax: let* bind pair requires identifier, but got %v", pair.Children[0].Type)), nil, nil
		}

		key := pair.Children[0].Str
		value := evalWithTailOptimization(pair.Children[1], env)

		env.define(key, value)
	}

	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, env)
	}
	return nil, sentences[len(sentences)-1], env
}

func evalCond(n *node.Node, env *env) (*object, *node.Node, *env) {
	if len(n.Children) == 1 {
		return newErrorObject("bad syntax: cond needs at least 1 argument, but got 0"), nil, nil
	}

	for _, branch := range n.Children[1:] {
		if branch.Type != node.Branch || len(branch.Children) == 0 {
			return newErrorObject("bad syntax: cond bad branch"), nil, nil
		}

		test := branch.Children[0]
		if (test.Type == node.Keyword && test.Str == "else") ||
			evalWithTailOptimization(test, env).isTruthy() {
			for _, child := range branch.Children[1 : len(branch.Children)-1] {
				evalWithTailOptimization(child, env)
			}
			return nil, branch.Children[len(branch.Children)-1], env
		}
	}
	// no cond match
	return voidObject, nil, nil
}

func evalSet(n *node.Node, env *env) *object {
	if len(n.Children) != 3 {
		return newErrorObject(fmt.Sprintf("set! exactly needs 2 arguments, but got %v", len(n.Children)-1))
	}
	if n.Children[1].Type != node.Identifier {
		return newErrorObject(fmt.Sprintf("1st argument of set! needs to be identifier, but got %v", n.Children[1].Type))
	}
	key := n.Children[1].Str
	value := evalWithTailOptimization(n.Children[2], env)
	if _, ok := env.lookup(key); !ok {
		return newErrorObject(fmt.Sprintf("set!: %v is not defined yet", key))
	}
	env.define(key, value)
	return voidObject
}

func evalQuote(n *node.Node) *object {
	switch n.Type {
	case node.Number:
		return newNumberObject(n.Num)
	case node.Boolean:
		return newBooleanObject(n.B)
	case node.String:
		return newStringObject(n.Str)
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

func evalDefine(n *node.Node, env *env) *object {
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
		return evalWithTailOptimization(&node.Node{Type: node.Branch, Children: []*node.Node{
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
	value := evalWithTailOptimization(n.Children[2], env)
	env.define(key, value)
	return voidObject
}

func evalLambda(n *node.Node, e *env) *object {
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
	return newRawFunctionObject(func(objects []*object) (*object, *node.Node, *env) {
		if len(objects) != len(argNames) {
			return newErrorObject(fmt.Sprintf("expected length of arguments to be %v, but got %v", len(argNames), len(objects))), nil, nil
		}

		newEnv := e.newEnv(newBindingFrame(argNames, objects))
		for _, sentence := range sentences[:len(sentences)-1] {
			evalWithTailOptimization(sentence, newEnv)
		}
		return nil, sentences[len(sentences)-1], newEnv
	})
}

func evalBegin(n *node.Node, env *env) (*object, *node.Node, *env) {
	if len(n.Children) <= 1 {
		return newErrorObject("begin needs at least 1 argument, but got 0"), nil, nil
	}

	sentences := n.Children[1:]
	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, env)
	}
	return nil, sentences[len(sentences)-1], env
}

// eval evaluates the given node, and returns the result obj, nil continuation, and nil newEnv.
// Otherwise, returns nil, continuation node, and newEnv to evaluate with for tail call optimization.
func eval(n *node.Node, env *env) (obj *object, continuation *node.Node, newEnv *env) {
	// Base cases
	switch n.Type {
	case node.Keyword:
		return newErrorObject("unexpected keyword"), nil, nil
	case node.Identifier:
		if obj, ok := env.lookup(n.Str); ok {
			return obj, nil, nil
		} else {
			return newErrorObject(fmt.Sprintf("unbound identifier: %v", n.Str)), nil, nil
		}
	case node.Number:
		return newNumberObject(n.Num), nil, nil
	case node.Boolean:
		return newBooleanObject(n.B), nil, nil
	case node.String:
		return newStringObject(n.Str), nil, nil
	}
	if n.Type != node.Branch {
		panic("node type not implemented")
	}

	// assert n.Type == node.Branch
	if len(n.Children) == 0 {
		return newErrorObject("bad syntax: empty sentence (\"()\") cannot be evaluated"), nil, nil
	}

	// Special forms
	if n.Children[0].Type == node.Keyword {
		switch n.Children[0].Str {
		case "and":
			return evalAnd(n, env), nil, nil
		case "or":
			return evalOr(n, env), nil, nil
		case "if":
			return evalIf(n, env)
		case "let":
			return evalLet(n, env)
		case "let*":
			return evalLetSeq(n, env)
		case "cond":
			return evalCond(n, env)
		case "set!":
			return evalSet(n, env), nil, nil
		case "quote":
			if len(n.Children) != 2 {
				return newErrorObject(fmt.Sprintf("quote needs exactly 1 argument, but got %v", len(n.Children)-1)), nil, nil
			}
			return evalQuote(n.Children[1]), nil, nil
		case "define":
			return evalDefine(n, env), nil, nil
		case "lambda":
			return evalLambda(n, env), nil, nil
		case "begin":
			// begin is not technically special form, but for tail optimization
			return evalBegin(n, env)
		}
	}

	// Function application
	objects := make([]*object, len(n.Children))
	for idx, child := range n.Children {
		objects[idx] = evalWithTailOptimization(child, env)
	}
	if objects[0].objectType != function {
		return newErrorObject(fmt.Sprintf("expected function in 0-th argument, but got %v", objects[0])), nil, nil
	}
	return objects[0].f(objects[1:])
}

func evalWithTailOptimization(n *node.Node, env *env) (ret *object) {
	for {
		ret, n, env = eval(n, env)
		if ret != nil {
			return
		}
	}
}
