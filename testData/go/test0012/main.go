//go:build test

package main

// Test when a type parameter has an approximate with function,
// 	 see [generics spec](https://go.dev/ref/spec#General_interfaces)

type X interface {
	~int
	String() string
}

type Y int

func (y Y) String() string {
	return `I'm an int!`
}

func Z[T X](x T) {
	println(x.String(), `=>`, int(x))
}

func main() {
	y := Y(42)
	Z(y)
}
