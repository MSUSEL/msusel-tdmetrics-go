package visitor

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	// Visitor is a handler for visiting a set of objects.
	Visitor struct {
		handle func(value any) bool
	}

	// Visitable is an object which has visitable children.
	Visitable interface {

		// VisitChildren will visit all off this object's children.
		//
		// This object will have already been visited and the visitor
		// will have already returned true for it.
		//
		// If visit returns false for any child then this should stop and
		// return false. Return true if visit returns true for all children.
		Visit(v Visitor) bool
	}
)

// New creates a new visitor for the given function handler.
func New(handle func(value any) bool) Visitor {
	return Visitor{handle: handle}
}

// Visit visits all the given values.
//
// If visit returns false for any value then this will stop and return false.
// It returns true if visit returns true for all values.
func Visit[T any](v Visitor, values ...T) bool {
	for _, value := range values {
		if !visitOne(v, value) {
			return false
		}
	}
	return true
}

// VisitList visits all the given values.
//
// If visit returns false for any value then this will stop and return false.
// It returns true if visit returns true for all values.
func VisitList[T any](v Visitor, values collections.ReadonlyList[T]) bool {
	for i := range values.Count() {
		if !visitOne(v, values.Get(i)) {
			return false
		}
	}
	return true
}

func visitOne[T any](v Visitor, value T) bool {
	if utils.IsNil(value) || !v.handle(value) {
		return false
	}
	if visitable, ok := any(value).(Visitable); ok {
		if !visitable.Visit(v) {
			return false
		}
	}
	return true
}
