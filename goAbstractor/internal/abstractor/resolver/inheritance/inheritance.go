package inheritance

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Node[T any] interface {
	// Implements determines if this interface implements the other interface.
	Implements(other T) bool

	AddInherits(parent T) T

	Inherits() collections.SortedSet[T]
}

type Inheritance[T Node[T]] interface {

	// Process adds a new node into the inheritance forest as a parent of
	// the correct nodes and having the correct parents added to it.
	//
	// This requires each node to be unique from any other node such that
	// any two different nodes will result in a non-zero from the given
	// comparer. This expects the node to have no parents already added
	// when passed into this method.
	Process(node T)
}

type inheritanceImp[T Node[T]] struct {
	roots collections.SortedSet[T]
	log   *logger.Logger
	count int
}

func New[T Node[T]](comp comp.Comparer[T], log *logger.Logger) Inheritance[T] {
	return &inheritanceImp[T]{
		roots: sortedSet.New(comp),
		log:   log,
	}
}

func (in *inheritanceImp[T]) Process(node T) {
	in.log.Logf(`╶─(%d) insert %v`, in.count, node)
	addParent[T](in.roots, node, in.log.Prefix(`  `))
	in.count++
}

func addParent[T Node[T]](siblings collections.SortedSet[T], n T, log *logger.Logger) {
	log2 := log.Prefix(` │ `)
	addedToSibling := false
	parentedSiblings := false
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)

		switch {
		case a.Implements(n):
			// Yi <: X, meaning `n` is a parent (sub-type) of `a`,
			// so add `n` as a parent of `a` and don't add it here.
			//
			// For example: {A, B} is a parent of {A, B, C} but {A, B, C}
			// may already have the parent {A} in it, so we have to recursively
			// call addParent to re-parent {A} as a parent of {A, B}.
			log.Logf(` ├─(%d) parent %v`, i, a)
			addParent(a.Inherits(), n, log2)
			addedToSibling = true

		case n.Implements(a):
			// Yi :> X, meaning `n` is a child (super-type) of `a`,
			// so move `a` from this set and add `a` as a parent of `n`.
			// This means that `n` is a parent in this set since otherwise
			// `a` would have been a parent to another object in this set.
			// We can simply add `a` to `n` since `a` has already been
			// checked against the other parents, hence it was in this set.
			log.Logf(` ├─(%d) child %v`, i, a)
			n.AddInherits(a)
			siblings.RemoveRange(i, 1)
			parentedSiblings = true

		default:
			// Possible overlap, check for parents (sub-types) in subtree.
			// Since we can't to overlaps in Go, just check any that aren't
			// specifically a super-type or sub-type.
			//
			// For example: {A, B, C} overlaps with {A, D} and {A, D} may
			// have the parents {A} and {D} in it. We want to add {A} as
			// a parent to {A, B, C}.
			log.Logf(` ├─(%d) else %v`, i, a)
			seekInherits(a.Inherits(), n, log2)
		}
	}

	if parentedSiblings {
		log.Log(` └─ add: parented sibling`)
		siblings.Add(n)
	} else if addedToSibling {
		log.Log(` └─ no-op: added to sibling`)
	} else {
		log.Log(` └─ add: default`)
		siblings.Add(n)
	}
}

func seekInherits[T Node[T]](siblings collections.SortedSet[T], n T, log *logger.Logger) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		if n.Implements(a) {
			log.Logf(` + %v`, a)
			n.AddInherits(a)
		} else {
			seekInherits(a.Inherits(), n, log)
		}
	}
}
