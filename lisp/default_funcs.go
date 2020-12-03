package lisp

import (
	"fmt"
	"github.com/motoki317/lisp-interpreter/node"
	"math"
	"strings"
)

var defaultEnv map[string]*object

func init() {
	defaultEnv = make(map[string]*object)

	defaultEnv["+"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			var res float64
			for _, in := range input {
				res += in
			}
			return newNumberObject(res)
		}))
	defaultEnv["-"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			if len(input) == 0 {
				return newErrorObject("expected at least one argument")
			}
			res := input[0]
			for _, in := range input[1:] {
				res -= in
			}
			return newNumberObject(res)
		}))
	defaultEnv["*"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			var res float64 = 1
			for _, in := range input {
				res *= in
			}
			return newNumberObject(res)
		}))
	defaultEnv["/"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			if len(input) == 0 {
				return newErrorObject("expected at least one argument")
			}
			res := input[0]
			for _, in := range input[1:] {
				if in == 0 {
					return newErrorObject("division by 0")
				}
				res /= in
			}
			return newNumberObject(res)
		}))

	defaultEnv[">"] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] > input[1])
		})))
	defaultEnv[">="] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] >= input[1])
		})))
	defaultEnv["="] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] == input[1])
		})))
	defaultEnv["<="] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] <= input[1])
		})))
	defaultEnv["<"] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] < input[1])
		})))

	defaultEnv["max"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			if len(input) == 0 {
				return newErrorObject("max: expected at least one input")
			}
			max := input[0]
			for _, in := range input[1:] {
				if max < in {
					max = in
				}
			}
			return newNumberObject(max)
		}))
	defaultEnv["min"] = newFunctionObject(
		makeNumbers(func(input []float64) *object {
			if len(input) == 0 {
				return newErrorObject("min: expected at least one input")
			}
			min := input[0]
			for _, in := range input[1:] {
				if in < min {
					min = in
				}
			}
			return newNumberObject(min)
		}))

	defaultEnv["zero?"] = newFunctionObject(
		makeUnary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] == 0)
		})))
	defaultEnv["even?"] = newFunctionObject(
		makeUnary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] == math.Trunc(input[0]) && int64(input[0])%2 == 0)
		})))
	defaultEnv["odd?"] = newFunctionObject(
		makeUnary(makeNumbers(func(input []float64) *object {
			return newBooleanObject(input[0] == math.Trunc(input[0]) && int64(input[0])%2 == 1)
		})))

	defaultEnv["modulo"] = newFunctionObject(
		makeBinary(makeNumbers(func(input []float64) *object {
			return newNumberObject(float64(int64(input[0]) % int64(input[1])))
		})))

	// and, or -> short circuit
	defaultEnv["not"] = newFunctionObject(
		makeUnary(makeBooleans(func(booleans []bool) *object {
			return newBooleanObject(!booleans[0])
		})))

	defaultEnv["sqrt"] = newFunctionObject(
		makeUnary(makeNumbers(func(input []float64) *object {
			return newNumberObject(math.Sqrt(input[0]))
		})))

	defaultEnv["cons"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			return newConsObject(objects[0], objects[1])
		}))
	defaultEnv["list"] = newFunctionObject(list)

	car := makeUnary(func(objects []*object) *object {
		o := objects[0]
		if o.objectType != cons {
			return newErrorObject(fmt.Sprintf("car: expected cons but got %v", o.objectType))
		}
		return o.pair[0]
	})
	cdr := makeUnary(func(objects []*object) *object {
		o := objects[0]
		if o.objectType != cons {
			return newErrorObject(fmt.Sprintf("cdr: expected cons but got %v", o))
		}
		return o.pair[1]
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
	defaultEnv["car"] = newFunctionObject(car)
	defaultEnv["cdr"] = newFunctionObject(cdr)
	second := crossProd(pair, pair)
	for k, v := range second {
		defaultEnv["c"+k+"r"] = newFunctionObject(composeFuncs(v...))
	}
	third := crossProd(second, pair)
	for k, v := range third {
		defaultEnv["c"+k+"r"] = newFunctionObject(composeFuncs(v...))
	}
	fourth := crossProd(third, pair)
	for k, v := range fourth {
		defaultEnv["c"+k+"r"] = newFunctionObject(composeFuncs(v...))
	}

	defaultEnv["set-car!"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			pair := objects[0]
			value := objects[1]
			if pair.objectType != cons {
				return newErrorObject(fmt.Sprintf("set-car!: expected 1st argument to be pair, but got %v", pair.objectType))
			}
			pair.pair[0] = value
			return voidObject
		}))
	defaultEnv["set-cdr!"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			pair := objects[0]
			value := objects[1]
			if pair.objectType != cons {
				return newErrorObject(fmt.Sprintf("set-cdr!: expected 1st argument to be pair, but got %v", pair.objectType))
			}
			pair.pair[1] = value
			return voidObject
		}))

	defaultEnv["equal?"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			o1, o2 := objects[0], objects[1]
			return newBooleanObject(o1.equals(o2))
		}))
	defaultEnv["eq?"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			o1, o2 := objects[0], objects[1]
			return newBooleanObject(o1.equals(o2))
		}))
	defaultEnv["eqv?"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			o1, o2 := objects[0], objects[1]
			return newBooleanObject(o1.equals(o2))
		}))
	defaultEnv["number?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0].objectType == number)
		}))
	defaultEnv["boolean?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0].objectType == boolean)
		}))
	defaultEnv["symbol?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0].objectType == symbol)
		}))
	defaultEnv["list?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0].isList())
		}))
	defaultEnv["null?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0] == nullObject)
		}))
	defaultEnv["string?"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			return newBooleanObject(objects[0].objectType == str)
		}))

	defaultEnv["apply"] = newRawFunctionObject(func(objects []*object) (*object, *node.Node, env) {
		if len(objects) != 2 {
			return newErrorObject(fmt.Sprintf("expected argument length to be 2, but got %v", len(objects))), nil, nil
		}
		f := objects[0]
		args := objects[1]
		if f.objectType != function {
			return newErrorObject(fmt.Sprintf("expected 1st argument of apply to be a function, but got %v", f)), nil, nil
		}
		if !args.isList() {
			return newErrorObject(fmt.Sprintf("expected 2nd argument of apply to be a list, but got %v", args)), nil, nil
		}
		return f.f(args.listElements())
	})
	defaultEnv["map"] = newFunctionObject(
		makeBinary(func(objects []*object) *object {
			f := objects[0]
			lst := objects[1]
			if f.objectType != function {
				return newErrorObject(fmt.Sprintf("expected 1st argument of map to be a function, but got %v", f))
			}
			if !lst.isList() {
				return newErrorObject(fmt.Sprintf("expected 2nd argument of map to be a list, but got %v", lst))
			}
			elements := lst.listElements()
			for i, elt := range elements {
				elements[i] = callWithTailOptimization(f.f, []*object{elt})
			}
			return list(elements)
		}))

	defaultEnv["symbol->string"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			o := objects[0]
			if o.objectType != symbol {
				return newErrorObject(fmt.Sprintf("expected 1st argument of symbol->string to be symbol, but got %v", o.objectType))
			}
			return newStringObject(o.str)
		}))
	defaultEnv["string->symbol"] = newFunctionObject(
		makeUnary(func(objects []*object) *object {
			o := objects[0]
			if o.objectType != str {
				return newErrorObject(fmt.Sprintf("expected 1st argument of string->symbol to be string, but got %v", o.objectType))
			}
			return newSymbolObject(o.str)
		}))
	defaultEnv["string-append"] = newFunctionObject(
		makeStrings(func(input []string) *object {
			return newStringObject(strings.Join(input, ""))
		}))
}

func list(objects []*object) *object {
	if len(objects) == 0 {
		return nullObject
	}
	return newConsObject(objects[0], list(objects[1:]))
}

// composeFuncs composes the given functions, applying from the LAST to the FIRST.
func composeFuncs(funcs ...generalFunc) generalFunc {
	return func(objects []*object) *object {
		for i := len(funcs) - 1; i >= 0; i-- {
			objects = []*object{funcs[i](objects)}
		}
		return objects[0]
	}
}

func makeNullary(next generalFunc) generalFunc {
	return func(objects []*object) *object {
		if len(objects) != 0 {
			return newErrorObject(fmt.Sprintf("expected length of argument to be 0, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeUnary(next generalFunc) generalFunc {
	return func(objects []*object) *object {
		if len(objects) != 1 {
			return newErrorObject(fmt.Sprintf("expected length of argument to be 1, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeBinary(next generalFunc) generalFunc {
	return func(objects []*object) *object {
		if len(objects) != 2 {
			return newErrorObject(fmt.Sprintf("expected length of argument to be 2, but got %v", len(objects)))
		}
		return next(objects)
	}
}

func makeNumbers(next func(input []float64) *object) generalFunc {
	return func(objects []*object) *object {
		nums := make([]float64, len(objects))
		for i, obj := range objects {
			if obj.objectType != number {
				return newErrorObject(fmt.Sprintf(
					"expected %v-th argument to be number, but got %v", i, obj))
			}
			nums[i] = obj.num
		}
		return next(nums)
	}
}

func makeBooleans(next func(input []bool) *object) generalFunc {
	return func(objects []*object) *object {
		booleans := make([]bool, len(objects))
		for i, obj := range objects {
			if obj.objectType != boolean {
				return newErrorObject(fmt.Sprintf(
					"expected %v-th argument to be boolean, but got %v", i, obj))
			}
			booleans[i] = obj.b
		}
		return next(booleans)
	}
}

func makeStrings(next func(input []string) *object) generalFunc {
	return func(objects []*object) *object {
		inputs := make([]string, len(objects))
		for i, obj := range objects {
			if obj.objectType != str {
				return newErrorObject(fmt.Sprintf(
					"expected %v-th argument to be str, but got %v", i, obj))
			}
			inputs[i] = obj.str
		}
		return next(inputs)
	}
}
