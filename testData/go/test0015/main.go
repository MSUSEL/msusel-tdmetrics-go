//go:build test

package main

// This is a subset of the test0011 to track down some issues
// in the type parameter and type arguments while creating instances.

func AsSlice[T comparable](t ...T) []T {
	return t
}

func main() {
	println(AsSlice(1, 2, 3, 4, 5))
}
