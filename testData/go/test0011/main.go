//go:build test

package main

// Test when a struct has a type param that is only used for another
// type param and an approximate type.
// Also more complex type parameters with pointers, arrays, and maps.

type Set[K comparable, V any, M ~map[K]*V] struct {
	m M
}

func (s Set[K, V, M]) AsSlices() ([]K, []*V) {
	keys := make([]K, len(s.m))
	values := make([]*V, len(s.m))
	i := 0
	for key, value := range s.m {
		keys[i] = key
		values[i] = value
		i++
	}
	return keys, values
}

func PrintSlice[T any, S ~[]T](s S) {
	print(`[ `)
	for i, v := range s {
		if i > 0 {
			print(`, `)
		}
		print(v)
	}
	println(` ]`)
}

type Bacon struct {
	Set[string, int, map[string]*int] // solidify types with embedded
}

func main() {
	x, y := 13, 15
	b := Bacon{
		Set: Set[string, int, map[string]*int]{
			m: map[string]*int{
				`hello`: &x,
				`world`: &y,
			},
		},
	}
	ks, vs := b.AsSlices()
	print(`keys:   `)
	PrintSlice(ks)
	print(`values: `)
	PrintSlice(vs)
}
