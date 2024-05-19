package typeDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"
)

// TypeDesc is an interface for all type descriptors.
type TypeDesc interface {

	// GoType gets the Go type associated with this type desc.
	GoType() types.Type

	// Equal indicates if the types are equivalent.
	Equal(other TypeDesc) bool
}

func equal[T TypeDesc](a, b T) bool {
	if m, n := utils.IsNil(a), utils.IsNil(b); m || n {
		return m && n
	}
	return a.Equal(b)
}

func equalTest[T TypeDesc](a T, b TypeDesc, h func(a, b T) bool) bool {
	if m, n := utils.IsNil(a), utils.IsNil(b); m || n {
		return m && n
	}
	b2, ok := b.(T)
	return ok && h(a, b2)
}

func equalList[T TypeDesc, S ~[]T](a, b S) bool {
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

func equalMap[T TypeDesc, M ~map[string]T](a, b M) bool {
	count := len(a)
	if count != len(b) {
		return false
	}
	for k, v := range a {
		if !v.Equal(b[k]) {
			return false
		}
	}
	return true
}
