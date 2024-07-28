//go:build test

package main

// Test when a struct has a type param used by a field
// that isn't used in any method.

type A[T int | float64 | string] struct {
	value T
}

func (a *A[T]) Mul(v int) {
	switch t := any(a.value).(type) {
	case int:
		a.value = any(t * v).(T)
	case float64:
		a.value = any(t * float64(v)).(T)
	case string:
		str := ``
		for range v {
			str += t
		}
		a.value = any(str).(T)
	default:
		panic(`unexpected type`)
	}
}

func main() {
	a := A[int]{value: 6}
	a.Mul(5)
	println(a.value) // 30

	b := A[float64]{value: 3.5}
	b.Mul(5)
	println(b.value) // +1.750000e+001

	c := A[string]{value: `Da`}
	c.Mul(5)
	println(c.value) // DaDaDaDaDa
}
