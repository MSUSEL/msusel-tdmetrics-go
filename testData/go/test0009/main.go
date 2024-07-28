//go:build test

package main

// Test when a struct has a type param that isn't used
// in any field but is used in a method on that struct.

type A[T int | float64 | string] struct {
	value int
}

func (a A[T]) Mul(v T) T {
	switch t := any(v).(type) {
	case int:
		return any(a.value * t).(T)
	case float64:
		return any(float64(a.value) * t).(T)
	case string:
		str := ``
		for range a.value {
			str += t
		}
		return any(str).(T)
	}
	panic(`unexpected type`)
}

func main() {
	println(A[int]{value: 5}.Mul(6))       // 30
	println(A[float64]{value: 5}.Mul(3.5)) // +1.750000e+001
	println(A[string]{value: 5}.Mul(`Da`)) // DaDaDaDaDa
}
