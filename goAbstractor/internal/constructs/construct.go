package constructs

import "github.com/Snow-Gremlin/goToolbox/utils"

type Visitor func(value Construct) bool

type Construct interface {
	Visit(v Visitor)

	// SetIndex sets the type index.
	SetIndex(index int)

	// Equal indicates if the types are equivalent.
	Equal(other Construct) bool
}

func visitTest[T Construct](v Visitor, value T) {
	if !utils.IsNil(value) && v(value) {
		value.Visit(v)
	}
}

func visitList[T Construct, S ~[]T](v Visitor, values S) {
	for _, val := range values {
		visitTest(v, val)
	}
}

func equal[T Construct](a, b T) bool {
	if m, n := utils.IsNil(a), utils.IsNil(b); m || n {
		return m && n
	}
	return a.Equal(b)
}

func equalTest[T Construct](a T, b Construct, h func(a, b T) bool) bool {
	if m, n := utils.IsNil(a), utils.IsNil(b); m || n {
		return m && n
	}
	b2, ok := b.(T)
	return ok && h(a, b2)
}

func equalList[T Construct, S ~[]T](a, b S) bool {
	count := len(a)
	if count != len(b) {
		return false
	}
	for i, v := range a {
		if !v.Equal(b[i]) {
			return false
		}
	}
	return true
}
