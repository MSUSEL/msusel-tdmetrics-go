//go:build test

package main

func dat[T any](x T, y int) {
	type nested[U any] struct {
		X T
		Y U
	}
	foo(nested[int]{X: x, Y: y})
}

func foo[T any](z T) {
	println(any(z))
}

func main() {
	dat(`four`, 2)
}
