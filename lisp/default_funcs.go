package lisp

import "fmt"

var defaultFuncs map[string]*object

func init() {
	defaultFuncs = make(map[string]*object)

	defaultFuncs["+"] = newFunctionObject(
		makeNumbersVariadicFunc(func(input []float64) *object {
			var res float64
			for _, in := range input {
				res += in
			}
			return newNumberObject(res)
		}))
	defaultFuncs["-"] = newFunctionObject(
		makeNumbersVariadicFunc(func(input []float64) *object {
			if len(input) == 0 {
				return newErrorObject("expected at least one argument")
			}
			res := input[0]
			for _, in := range input[1:] {
				res -= in
			}
			return newNumberObject(res)
		}))
	defaultFuncs["*"] = newFunctionObject(
		makeNumbersVariadicFunc(func(input []float64) *object {
			var res float64 = 1
			for _, in := range input {
				res *= in
			}
			return newNumberObject(res)
		}))
	defaultFuncs["/"] = newFunctionObject(
		makeNumbersVariadicFunc(func(input []float64) *object {
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
}

func makeNumbersVariadicFunc(next func(input []float64) *object) generalFunc {
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
