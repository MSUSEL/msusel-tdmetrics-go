//go:build test

package main

func sum(values ...int) int {
	s := 0
	for _, v := range values {
		s += v
	}
	return s
}

func first(values ...int) int {
	return values[0]
}

func last(values ...int) int {
	return values[len(values)-1]
}

func main() {
	data := []int{32, 54, 8, 133, 75}
	print(`sum: `, sum(data...))
	print(`first: `, first(data...))
	print(`last: `, last(data...))
}
