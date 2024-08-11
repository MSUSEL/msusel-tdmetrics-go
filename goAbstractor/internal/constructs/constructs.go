package constructs

import "github.com/Snow-Gremlin/goToolbox/comp"

// Construct is a description of a type.
type Construct interface {
	comp.Comparable[Construct]

	// Kind gets a string unique to each construct type.
	Kind() string

	// setIndex sets the unique index of construct.
	// Indices will be 1 based so that 0 is unset.
	SetIndex(index int)
}

func Comparer[T Construct]() comp.Comparer[T] {
	cmp := comp.ComparableComparer[Construct]()
	return func(x, y T) int { return cmp(x, y) }
}

func ComparerPend[T Construct](a, b T) func() int {
	return Comparer[T]().Pend(a, b)
}

func SliceComparer[T Construct]() comp.Comparer[[]T] {
	return comp.Slice[[]T](Comparer[T]())
}

func SliceComparerPend[T Construct](a, b []T) func() int {
	return SliceComparer[T]().Pend(a, b)
}

func CompareTo[T Construct](a T, b Construct, cmp comp.Comparer[T]) int {
	return comp.Or(
		comp.DefaultPend(a.Kind(), b.Kind()),
		cmp.Pend(a, b.(T)),
	)
}
