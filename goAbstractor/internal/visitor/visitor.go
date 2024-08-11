package visitor

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Visitor is a handler for visiting a set of objects.
type Visitor struct {
	handle func(value any) bool
}

// Visitable is an object which has visitable children.
type Visitable interface {

	// VisitChildren will visit all off this object's children.
	//
	// This object will have already been visited and the visitor
	// will have already returned true for it.
	Visit(v Visitor)
}

// New creates a new visitor for the given function handler.
func New(handle func(value any) bool) Visitor {
	return Visitor{handle: handle}
}

// Visit visits all the given values.
func Visit[T any](v Visitor, values ...T) {
	for _, value := range values {
		visitOne(v, value)
	}
}

// VisitList visits all the given values.
func VisitList[T any](v Visitor, values collections.ReadonlyList[T]) {
	for i := range values.Count() {
		visitOne(v, values.Get(i))
	}
}

func visitOne[T any](v Visitor, value T) {
	if !utils.IsNil(value) && v.handle(value) {
		if visitable, ok := any(value).(Visitable); ok {
			visitable.Visit(v)
		}
	}
}
