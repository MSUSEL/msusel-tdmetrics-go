//go:build test

package main

func Len[T any](start, stop *Node[T]) int {
	count := 0
	for n := start; n != stop; n = n.next {
		count++
	}
	return count
}

type Node[T any] struct {
	next  *Node[T]
	value T
}

func main() {
	s := &Node[int]{value: 1, next: &Node[int]{value: 2, next: &Node[int]{value: 3}}}
	println(Len(s, nil))
}
