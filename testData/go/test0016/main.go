//go:build test

package main

type Node[T any] interface {
	comparable
	Next() T
}

func Len[T Node[T]](start, stop T) int {
	count := 0
	for n := start; n != stop; n = n.Next() {
		count++
	}
	return count
}

type nodeImp struct {
	next *nodeImp
}

func (ni *nodeImp) Next() *nodeImp {
	return ni
}

func main() {
	s := &nodeImp{next: &nodeImp{next: &nodeImp{}}}
	println(Len(s, nil))
}
