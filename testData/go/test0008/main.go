//go:build test

package main

type Foo[T any] struct {
	value T
}

func (f *Foo[X]) Get() X { // X should be replaced with T
	return f.value
}

func main() {
	f := &Foo[int]{value: 42}
	print(`value: `, f.Get())
}
