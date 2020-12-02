package lisp

import (
	"fmt"
	"math"
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

	// and, or -> short circuit
	defaultEnv["not"] = newFunctionObject(
		makeUnary(makeBooleans(func(booleans []bool) *object {
			return newBooleanObject(!booleans[0])
		})))

	defaultEnv["sqrt"] = newFunctionObject(
		makeUnary(makeNumbers(func(input []float64) *object {
			return newNumberObject(math.Sqrt(input[0]))
		})))
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
