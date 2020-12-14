package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/lisp/macro"
	"github.com/motoki317/lisp-interpreter/lisp/object"
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
	"time"
)

func evalAnd(n *node.Node, env *object.Env) object.Object {
	// Short circuit evaluation
	res := object.NewBooleanObject(true)
	for _, child := range n.Children[1:] {
		res = evalWithTailOptimization(child, env)
		if !res.IsTruthy() {
			return object.NewBooleanObject(false)
		}
	}
	return res
}

func evalOr(n *node.Node, env *object.Env) object.Object {
	// Short circuit evaluation
	for _, child := range n.Children[1:] {
		res := evalWithTailOptimization(child, env)
		if res.IsTruthy() {
			return res
		}
	}
	return object.NewBooleanObject(false)
}

func evalIf(n *node.Node, env *object.Env) (object.Object, *node.Node, *object.Env) {
	if len(n.Children) != 3 && len(n.Children) != 4 {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: if needs 2 or 3 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	res := evalWithTailOptimization(n.Children[1], env)
	if res.IsTruthy() {
		return nil, n.Children[2], env
	} else {
		if len(n.Children) == 4 {
			return nil, n.Children[3], env
		} else {
			return object.VoidObj, nil, nil
		}
	}
}

func evalLet(n *node.Node, e *object.Env) (object.Object, *node.Node, *object.Env) {
	if len(n.Children) <= 2 {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: let needs at least 2 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	keys := make([]string, len(pairs.Children))
	values := make([]object.Object, len(pairs.Children))
	for i, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return object.NewErrorObject(fmt.Sprintf("bad syntax: let bind pair needs a list of length 2, but got length %v", len(pair.Children))), nil, nil
		}
		if pair.Children[0].Type != node.Identifier {
			return object.NewErrorObject(fmt.Sprintf("bad syntax: let bind pair requires identifier, but got %v", pair.Children[0].Type)), nil, nil
		}

		keys[i] = pair.Children[0].Str
		values[i] = evalWithTailOptimization(pair.Children[1], e)
	}

	e = e.NewEnv(object.NewBindingFrame(keys, values))
	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, e)
	}
	return nil, sentences[len(sentences)-1], e
}

func evalLetSeq(n *node.Node, e *object.Env) (object.Object, *node.Node, *object.Env) {
	if len(n.Children) <= 2 {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: let* needs at least 2 arguments, but got %v", len(n.Children)-1)), nil, nil
	}

	pairs := n.Children[1]
	sentences := n.Children[2:]

	e = e.NewEnv(object.EmptyFrame())
	for _, pair := range pairs.Children {
		if len(pair.Children) != 2 {
			return object.NewErrorObject(fmt.Sprintf("bad syntax: let* bind pair needs a list of length 2, but got length %v", len(pair.Children))), nil, nil
		}
		if pair.Children[0].Type != node.Identifier {
			return object.NewErrorObject(fmt.Sprintf("bad syntax: let* bind pair requires identifier, but got %v", pair.Children[0].Type)), nil, nil
		}

		key := pair.Children[0].Str
		value := evalWithTailOptimization(pair.Children[1], e)

		e.Define(key, value)
	}

	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, e)
	}
	return nil, sentences[len(sentences)-1], e
}

func evalCond(n *node.Node, env *object.Env) (object.Object, *node.Node, *object.Env) {
	if len(n.Children) == 1 {
		return object.NewErrorObject("bad syntax: cond needs at least 1 argument, but got 0"), nil, nil
	}

	for _, branch := range n.Children[1:] {
		if branch.Type != node.Branch || len(branch.Children) == 0 {
			return object.NewErrorObject("bad syntax: cond bad branch"), nil, nil
		}

		test := branch.Children[0]
		if (test.Type == node.Keyword && test.Str == "else") ||
			evalWithTailOptimization(test, env).IsTruthy() {
			for _, child := range branch.Children[1 : len(branch.Children)-1] {
				evalWithTailOptimization(child, env)
			}
			return nil, branch.Children[len(branch.Children)-1], env
		}
	}
	// no cond match
	return object.VoidObj, nil, nil
}

func evalSet(n *node.Node, e *object.Env) object.Object {
	if len(n.Children) != 3 {
		return object.NewErrorObject(fmt.Sprintf("set! exactly needs 2 arguments, but got %v", len(n.Children)-1))
	}
	if n.Children[1].Type != node.Identifier {
		return object.NewErrorObject(fmt.Sprintf("1st argument of set! needs to be identifier, but got %v", n.Children[1].Type))
	}
	key := n.Children[1].Str
	value := evalWithTailOptimization(n.Children[2], e)
	if ok := e.Set(key, value); !ok {
		return object.NewErrorObject(fmt.Sprintf("set!: %v is not defined yet", key))
	}
	return object.VoidObj
}

func evalQuote(n *node.Node) object.Object {
	switch n.Type {
	case node.Number:
		return object.NewNumberObject(n.Num)
	case node.Boolean:
		return object.NewBooleanObject(n.B)
	case node.String:
		return object.NewStringObject(n.Str)
	case node.Identifier:
		return object.NewSymbolObject(n.Str)
	case node.Keyword:
		return object.NewSymbolObject(n.Str)
	}
	if n.Type != node.Branch {
		panic(fmt.Sprintf("quote node type not implemented: %v", n.Type))
	}
	// assert n.Type == node.Branch
	if len(n.Children) == 0 {
		return object.NullObj
	}
	if len(n.Children) == 3 && n.Children[1].Type == node.Keyword && n.Children[1].Str == "." {
		return object.NewConsObject(
			evalQuote(n.Children[0]),
			evalQuote(n.Children[2]))
	}
	return object.NewConsObject(
		evalQuote(n.Children[0]),
		evalQuote(&node.Node{
			Type:     node.Branch,
			Children: n.Children[1:],
		}))
}

func evalDefine(n *node.Node, e *object.Env) object.Object {
	// define syntax sugar
	// (define (func-name arg1 arg2) ...)
	// = (define func-name (lambda (arg1 arg2) ...))
	if n.Children[1].Type == node.Branch {
		if len(n.Children[1].Children) == 0 || n.Children[1].Children[0].Type != node.Identifier {
			return object.NewErrorObject("bad syntax: function definition requires function name")
		}

		funcName := n.Children[1].Children[0]
		argNames := n.Children[1].Children[1:]
		sentences := n.Children[2:]

		// specific define syntax sugar for variadic arguments function
		// (define (func-name . x) ...)
		// = (define func-name (lambda x ...))
		if len(argNames) == 2 &&
			argNames[0].Type == node.Keyword && argNames[0].Str == "." &&
			argNames[1].Type == node.Identifier {
			lambda := append([]*node.Node{
				{Type: node.Keyword, Str: "lambda"},
				argNames[1],
			}, sentences...)
			return evalDefine(&node.Node{Type: node.Branch, Children: []*node.Node{
				{Type: node.Keyword, Str: "define"},
				funcName,
				{Type: node.Branch, Children: lambda},
			}}, e)
		}

		// Rewrite AST and eval again
		lambda := append([]*node.Node{
			{Type: node.Keyword, Str: "lambda"},
			{Type: node.Branch, Children: argNames},
		}, sentences...)
		return evalDefine(&node.Node{Type: node.Branch, Children: []*node.Node{
			{Type: node.Keyword, Str: "define"},
			funcName,
			{Type: node.Branch, Children: lambda},
		}}, e)
	}

	// Normal define
	if len(n.Children) != 3 {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: define takes exactly 2 arguments, but got %v", len(n.Children)-1))
	}

	if n.Children[1].Type != node.Identifier {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: expected 1st argument of define to be identifier, but got %v", n.Children[1]))
	}

	key := n.Children[1].Str
	value := evalWithTailOptimization(n.Children[2], e)
	e.Define(key, value)
	return object.VoidObj
}

func evalLambda(n *node.Node, e *object.Env) object.Object {
	if len(n.Children) < 3 {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: lambda takes 2 or more arguments, but got %v", len(n.Children)-1))
	}

	sentences := n.Children[2:]

	// Variadic length arguments
	// (lambda x ...)
	if n.Children[1].Type == node.Identifier {
		lstName := n.Children[1].Str
		return object.NewFunctionObject(func(objects []object.Object) (object.Object, *node.Node, *object.Env) {
			newEnv := e.NewEnv(object.EmptyFrame())
			newEnv.Define(lstName, list(objects))
			for _, sentence := range sentences[:len(sentences)-1] {
				evalWithTailOptimization(sentence, newEnv)
			}
			return nil, sentences[len(sentences)-1], newEnv
		})
	}

	if n.Children[1].Type != node.Branch {
		return object.NewErrorObject(fmt.Sprintf("bad syntax: 1st argument of lambda needs to be a list of arguments, but got %v", n.Children[1]))
	}
	inputArgs := n.Children[1].Children

	// Variadic length arguments with leading arguments
	// (lambda (x y . z) ...)
	if len(inputArgs) >= 3 &&
		inputArgs[len(inputArgs)-2].Type == node.Keyword &&
		inputArgs[len(inputArgs)-2].Str == "." &&
		inputArgs[len(inputArgs)-1].Type == node.Identifier {
		argNames := make([]string, len(inputArgs)-2)
		for i := range argNames {
			if inputArgs[i].Type != node.Identifier {
				return object.NewErrorObject(fmt.Sprintf("bad syntax: expected %v-th argument of lambda function to be identifier, but got %v", i, inputArgs[i].Type))
			}
			argNames[i] = inputArgs[i].Str
		}
		lstName := inputArgs[len(inputArgs)-1].Str
		return object.NewFunctionObject(func(objects []object.Object) (object.Object, *node.Node, *object.Env) {
			if len(objects) < len(argNames) {
				return object.NewErrorObject(fmt.Sprintf("expected length of arguments to be greater than or equal to %v, but got %v", len(argNames), len(objects))), nil, nil
			}
			newEnv := e.NewEnv(object.NewBindingFrame(argNames, objects[:len(argNames)]))
			newEnv.Define(lstName, list(objects[len(argNames):]))
			for _, sentence := range sentences[:len(sentences)-1] {
				evalWithTailOptimization(sentence, newEnv)
			}
			return nil, sentences[len(sentences)-1], newEnv
		})
	}

	// Normal define
	// (lambda (x y z) ...)
	argNames := make([]string, len(inputArgs))
	for i, arg := range inputArgs {
		if arg.Type != node.Identifier {
			return object.NewErrorObject(fmt.Sprintf("bad syntax: expected %v-th argument of lambda function to be identifier, but got %v", i, arg.Type))
		}
		argNames[i] = arg.Str
	}
	return object.NewFunctionObject(func(objects []object.Object) (object.Object, *node.Node, *object.Env) {
		if len(objects) != len(argNames) {
			return object.NewErrorObject(fmt.Sprintf("expected length of arguments to be %v, but got %v", len(argNames), len(objects))), nil, nil
		}

		newEnv := e.NewEnv(object.NewBindingFrame(argNames, objects))
		for _, sentence := range sentences[:len(sentences)-1] {
			evalWithTailOptimization(sentence, newEnv)
		}
		return nil, sentences[len(sentences)-1], newEnv
	})
}

func evalBegin(n *node.Node, env *object.Env) (object.Object, *node.Node, *object.Env) {
	if len(n.Children) <= 1 {
		return object.NewErrorObject("begin needs at least 1 argument, but got 0"), nil, nil
	}

	sentences := n.Children[1:]
	for _, sentence := range sentences[:len(sentences)-1] {
		evalWithTailOptimization(sentence, env)
	}
	return nil, sentences[len(sentences)-1], env
}

func evalMacro(n *node.Node, e *object.Env) object.Object {
	m, err := macro.NewMacro(n)
	if err != nil {
		return object.NewErrorObject("bad macro syntax: " + err.Error())
	}
	e.DefineGlobalMacro(m)
	return object.VoidObj
}

func evalDelay(n *node.Node, e *object.Env) object.Object {
	if len(n.Children) != 2 {
		return object.NewErrorObject(fmt.Sprintf("delay needs exactly 1 argument, but got %v", len(n.Children)-1))
	}
	toDelay := n.Children[1]
	return object.NewPromiseObject(toDelay, e)
}

// eval evaluates the given node, and returns the result obj, nil continuation, and nil newEnv.
// Otherwise, returns nil, continuation node, and newEnv to evaluate with for tail call optimization.
func eval(n *node.Node, e *object.Env) (obj object.Object, continuation *node.Node, newEnv *object.Env) {
	// Base cases
	switch n.Type {
	case node.Keyword:
		return object.NewErrorObject("unexpected keyword"), nil, nil
	case node.Identifier:
		if obj, ok := e.Lookup(n.Str); ok {
			return obj, nil, nil
		} else {
			return object.NewErrorObject(fmt.Sprintf("unbound identifier: %v", n.Str)), nil, nil
		}
	case node.Number:
		return object.NewNumberObject(n.Num), nil, nil
	case node.Boolean:
		return object.NewBooleanObject(n.B), nil, nil
	case node.String:
		return object.NewStringObject(n.Str), nil, nil
	}
	if n.Type != node.Branch {
		panic("node type not implemented")
	}

	// assert n.Type == node.Branch
	if len(n.Children) == 0 {
		return object.NewErrorObject("bad syntax: empty sentence (\"()\") cannot be evaluated"), nil, nil
	}

	// Special forms
	if n.Children[0].Type == node.Keyword {
		switch n.Children[0].Str {
		case "and":
			return evalAnd(n, e), nil, nil
		case "or":
			return evalOr(n, e), nil, nil
		case "if":
			return evalIf(n, e)
		case "let":
			return evalLet(n, e)
		case "let*":
			return evalLetSeq(n, e)
		case "cond":
			return evalCond(n, e)
		case "set!":
			return evalSet(n, e), nil, nil
		case "quote":
			if len(n.Children) != 2 {
				return object.NewErrorObject(fmt.Sprintf("quote needs exactly 1 argument, but got %v", len(n.Children)-1)), nil, nil
			}
			return evalQuote(n.Children[1]), nil, nil
		case "define":
			return evalDefine(n, e), nil, nil
		case "lambda":
			return evalLambda(n, e), nil, nil
		case "begin":
			// begin is not technically special form, but for tail optimization
			return evalBegin(n, e)
		case "define-syntax":
			return evalMacro(n, e), nil, nil
		case "delay":
			return evalDelay(n, e), nil, nil
		}
	}

	// Function application
	objects := make([]object.Object, len(n.Children))
	for idx, child := range n.Children {
		objects[idx] = evalWithTailOptimization(child, e)
	}
	if objects[0].Type() != object_type.Function {
		return object.NewErrorObject(fmt.Sprintf("expected function in 0-th argument, but got %v", objects[0])), nil, nil
	}
	return objects[0].F(objects[1:])
}

func evalWithTailOptimization(n *node.Node, env *object.Env) (ret object.Object) {
	for {
		ret, n, env = eval(n, env)
		if ret != nil {
			return
		}
	}
}

func evalWithStopper(n *node.Node, env *object.Env, stop <-chan time.Time) (ret object.Object, timedOut bool) {
	for {
		select {
		case <-stop:
			return nil, true
		default:
			ret, n, env = eval(n, env)
			if ret != nil {
				return ret, false
			}
		}
	}
}
