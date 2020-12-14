package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/lisp/object"
	"github.com/motoki317/lisp-interpreter/lisp/object/object_type"
	"github.com/motoki317/lisp-interpreter/node"
	"math"
	"strings"
)

var defaultEnv map[string]object.Object

func init() {
	defaultEnv = make(map[string]object.Object)

	defaultEnv["+"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			var res float64
			for _, in := range input {
				res += in
			}
			return object.NewNumberObject(res)
		}))
	defaultEnv["-"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			if len(input) == 0 {
				return object.NewErrorObject("expected at least one argument")
			}
			res := input[0]
			for _, in := range input[1:] {
				res -= in
			}
			return object.NewNumberObject(res)
		}))
	defaultEnv["*"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			var res float64 = 1
			for _, in := range input {
				res *= in
			}
			return object.NewNumberObject(res)
		}))
	defaultEnv["/"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			if len(input) == 0 {
				return object.NewErrorObject("expected at least one argument")
			}
			res := input[0]
			for _, in := range input[1:] {
				if in == 0 {
					return object.NewErrorObject("division by 0")
				}
				res /= in
			}
			return object.NewNumberObject(res)
		}))

	defaultEnv[">"] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] > input[1])
		})))
	defaultEnv[">="] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] >= input[1])
		})))
	defaultEnv["="] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] == input[1])
		})))
	defaultEnv["<="] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] <= input[1])
		})))
	defaultEnv["<"] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] < input[1])
		})))

	defaultEnv["max"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			if len(input) == 0 {
				return object.NewErrorObject("max: expected at least one input")
			}
			max := input[0]
			for _, in := range input[1:] {
				if max < in {
					max = in
				}
			}
			return object.NewNumberObject(max)
		}))
	defaultEnv["min"] = object.NewWrappedFunctionObject(
		makeNumbers(func(input []float64) object.Object {
			if len(input) == 0 {
				return object.NewErrorObject("min: expected at least one input")
			}
			min := input[0]
			for _, in := range input[1:] {
				if in < min {
					min = in
				}
			}
			return object.NewNumberObject(min)
		}))

	defaultEnv["zero?"] = object.NewWrappedFunctionObject(
		makeUnary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] == 0)
		})))
	defaultEnv["even?"] = object.NewWrappedFunctionObject(
		makeUnary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] == math.Trunc(input[0]) && int64(input[0])%2 == 0)
		})))
	defaultEnv["odd?"] = object.NewWrappedFunctionObject(
		makeUnary(makeNumbers(func(input []float64) object.Object {
			return object.NewBooleanObject(input[0] == math.Trunc(input[0]) && int64(input[0])%2 == 1)
		})))

	defaultEnv["modulo"] = object.NewWrappedFunctionObject(
		makeBinary(makeNumbers(func(input []float64) object.Object {
			return object.NewNumberObject(float64(int64(input[0]) % int64(input[1])))
		})))

	// and, or -> short circuit
	defaultEnv["not"] = object.NewWrappedFunctionObject(
		makeUnary(makeBooleans(func(booleans []bool) object.Object {
			return object.NewBooleanObject(!booleans[0])
		})))

	defaultEnv["sqrt"] = object.NewWrappedFunctionObject(
		makeUnary(makeNumbers(func(input []float64) object.Object {
			return object.NewNumberObject(math.Sqrt(input[0]))
		})))

	defaultEnv["cons"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			return object.NewConsObject(objects[0], objects[1])
		}))
	defaultEnv["list"] = object.NewWrappedFunctionObject(list)

	car := makeUnary(func(objects []object.Object) object.Object {
		o := objects[0]
		if o.Type() != object_type.Cons {
			return object.NewErrorObject(fmt.Sprintf("car: expected cons but got %v", o.Type()))
		}
		return o.Pair()[0]
	})
	cdr := makeUnary(func(objects []object.Object) object.Object {
		o := objects[0]
		if o.Type() != object_type.Cons {
			return object.NewErrorObject(fmt.Sprintf("cdr: expected cons but got %v", o))
		}
		return o.Pair()[1]
	})

	pair := map[string][]generalFunc{"a": {car}, "d": {cdr}}
	crossProd := func(first, second map[string][]generalFunc) (res map[string][]generalFunc) {
		res = make(map[string][]generalFunc, len(first)*len(second))
		for k1, v1 := range first {
			for k2, v2 := range second {
				v := make([]generalFunc, 0, len(v1)+len(v2))
				res[k1+k2] = append(append(v, v1...), v2...)
			}
		}
		return
	}
	defaultEnv["car"] = object.NewWrappedFunctionObject(car)
	defaultEnv["cdr"] = object.NewWrappedFunctionObject(cdr)
	second := crossProd(pair, pair)
	for k, v := range second {
		defaultEnv["c"+k+"r"] = object.NewWrappedFunctionObject(composeFuncs(v...))
	}
	third := crossProd(second, pair)
	for k, v := range third {
		defaultEnv["c"+k+"r"] = object.NewWrappedFunctionObject(composeFuncs(v...))
	}
	fourth := crossProd(third, pair)
	for k, v := range fourth {
		defaultEnv["c"+k+"r"] = object.NewWrappedFunctionObject(composeFuncs(v...))
	}

	defaultEnv["set-car!"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			pair := objects[0]
			value := objects[1]
			if pair.Type() != object_type.Cons {
				return object.NewErrorObject(fmt.Sprintf("set-car!: expected 1st argument to be pair, but got %v", pair.Type()))
			}
			pair.Pair()[0] = value
			return object.VoidObj
		}))
	defaultEnv["set-cdr!"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			pair := objects[0]
			value := objects[1]
			if pair.Type() != object_type.Cons {
				return object.NewErrorObject(fmt.Sprintf("set-cdr!: expected 1st argument to be pair, but got %v", pair.Type()))
			}
			pair.Pair()[1] = value
			return object.VoidObj
		}))

	defaultEnv["equal?"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			o1, o2 := objects[0], objects[1]
			return object.NewBooleanObject(o1.Equals(o2))
		}))
	defaultEnv["eq?"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			o1, o2 := objects[0], objects[1]
			return object.NewBooleanObject(o1.Equals(o2))
		}))
	defaultEnv["eqv?"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			o1, o2 := objects[0], objects[1]
			return object.NewBooleanObject(o1.Equals(o2))
		}))
	defaultEnv["number?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].Type() == object_type.Number)
		}))
	defaultEnv["boolean?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].Type() == object_type.Boolean)
		}))
	defaultEnv["symbol?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].Type() == object_type.Symbol)
		}))
	defaultEnv["list?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].IsList())
		}))
	defaultEnv["null?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].Type() == object_type.Null)
		}))
	defaultEnv["string?"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			return object.NewBooleanObject(objects[0].Type() == object_type.Str)
		}))

	defaultEnv["apply"] = object.NewFunctionObject(func(objects []object.Object) (object.Object, *node.Node, *object.Env) {
		if len(objects) != 2 {
			return object.NewErrorObject(fmt.Sprintf("expected argument length to be 2, but got %v", len(objects))), nil, nil
		}
		f := objects[0]
		args := objects[1]
		if f.Type() != object_type.Function {
			return object.NewErrorObject(fmt.Sprintf("expected 1st argument of apply to be a function, but got %v", f)), nil, nil
		}
		if !args.IsList() {
			return object.NewErrorObject(fmt.Sprintf("expected 2nd argument of apply to be a list, but got %v", args)), nil, nil
		}
		return f.F(args.ListElements())
	})
	defaultEnv["map"] = object.NewWrappedFunctionObject(
		makeBinary(func(objects []object.Object) object.Object {
			f := objects[0]
			lst := objects[1]
			if f.Type() != object_type.Function {
				return object.NewErrorObject(fmt.Sprintf("expected 1st argument of map to be a function, but got %v", f))
			}
			if !lst.IsList() {
				return object.NewErrorObject(fmt.Sprintf("expected 2nd argument of map to be a list, but got %v", lst))
			}
			elements := lst.ListElements()
			for i, elt := range elements {
				elements[i] = callWithTailOptimization(f.F, []object.Object{elt})
			}
			return list(elements)
		}))

	defaultEnv["force"] = object.NewFunctionObject(func(objects []object.Object) (object.Object, *node.Node, *object.Env) {
		if len(objects) != 1 {
			return object.NewErrorObject(fmt.Sprintf("force needs exactly 1 argument, but got %v", len(objects))), nil, nil
		}
		o := objects[0]
		if o.Type() != object_type.Promise {
			return object.NewErrorObject(fmt.Sprintf("force takes promise object as argument, but got %v", o.Type())), nil, nil
		}
		return o.F(nil)
	})

	defaultEnv["symbol->string"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			o := objects[0]
			if o.Type() != object_type.Symbol {
				return object.NewErrorObject(fmt.Sprintf("expected 1st argument of symbol->string to be symbol, but got %v", o.Type()))
			}
			return object.NewStringObject(o.Str())
		}))
	defaultEnv["string->symbol"] = object.NewWrappedFunctionObject(
		makeUnary(func(objects []object.Object) object.Object {
			o := objects[0]
			if o.Type() != object_type.Str {
				return object.NewErrorObject(fmt.Sprintf("expected 1st argument of string->symbol to be string, but got %v", o.Type()))
			}
			return object.NewSymbolObject(o.Str())
		}))
	defaultEnv["string-append"] = object.NewWrappedFunctionObject(
		makeStrings(func(input []string) object.Object {
			return object.NewStringObject(strings.Join(input, ""))
		}))
}

type generalFunc func(objects []object.Object) object.Object

func callWithTailOptimization(f func(objects []object.Object) (object.Object, *node.Node, *object.Env), objects []object.Object) object.Object {
	obj, n, e := f(objects)
	if obj != nil {
		return obj
	}
	return evalWithTailOptimization(n, e)
}

func list(objects []object.Object) object.Object {
	if len(objects) == 0 {
		return object.NullObj
	}
	return object.NewConsObject(objects[0], list(objects[1:]))
}

// composeFuncs composes the given functions, applying from the LAST to the FIRST.
func composeFuncs(funcs ...generalFunc) generalFunc {
	return func(objects []object.Object) object.Object {
		for i := len(funcs) - 1; i >= 0; i-- {
			objects = []object.Object{funcs[i](objects)}
		}
		return objects[0]
	}
}

func makeNullary(next generalFunc) generalFunc {
	return func(objects []object.Object) object.Object {
		if len(objects) != 0 {
			return object.NewErrorObject(fmt.Sprintf("expected length of argument to be 0, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeUnary(next generalFunc) generalFunc {
	return func(objects []object.Object) object.Object {
		if len(objects) != 1 {
			return object.NewErrorObject(fmt.Sprintf("expected length of argument to be 1, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeBinary(next generalFunc) generalFunc {
	return func(objects []object.Object) object.Object {
		if len(objects) != 2 {
			return object.NewErrorObject(fmt.Sprintf("expected length of argument to be 2, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeNumbers(next func(input []float64) object.Object) generalFunc {
	return func(objects []object.Object) object.Object {
		nums := make([]float64, len(objects))
		for i, obj := range objects {
			if obj.Type() != object_type.Number {
				return object.NewErrorObject(fmt.Sprintf(
					"expected %v-th argument to be number, but got %v", i, obj))
			}
			nums[i] = obj.Number()
		}
		return next(nums)
	}
}

func makeBooleans(next func(input []bool) object.Object) generalFunc {
	return func(objects []object.Object) object.Object {
		booleans := make([]bool, len(objects))
		for i, obj := range objects {
			if obj.Type() != object_type.Boolean {
				return object.NewErrorObject(fmt.Sprintf(
					"expected %v-th argument to be boolean, but got %v", i, obj))
			}
			booleans[i] = obj.Bool()
		}
		return next(booleans)
	}
}

func makeStrings(next func(input []string) object.Object) generalFunc {
	return func(objects []object.Object) object.Object {
		inputs := make([]string, len(objects))
		for i, obj := range objects {
			if obj.Type() != object_type.Str {
				return object.NewErrorObject(fmt.Sprintf(
					"expected %v-th argument to be str, but got %v", i, obj))
			}
			inputs[i] = obj.Str()
		}
		return next(inputs)
	}
}
