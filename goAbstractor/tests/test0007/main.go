//go:build test

package main

type Foo[T string | int | uint] struct {
	value T
}

func (f *Foo[T]) Add(v T) T {
	f.value += v
	return f.value
}

func New[T string | int](v T) *Foo[T] {
	return &Foo[T]{value: v}
}

func main() {
	c := New(`Hello`)
	c.Add(`World`)
}
