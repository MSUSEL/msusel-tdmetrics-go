package inheritance

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

func Resolve[T Node[T]](log *logger.Logger, cmp comp.Comparer[T], its collections.ReadonlySortedSet[T]) {
	log2 := log.Group(`inheritance`).Indent()
	in := New(cmp, log2)
	for i := range its.Count() {
		in.Process(its.Get(i))
	}
	log2.Log()
}

type relationship int

const (
	// unknown indicates that a relationship hasn't been determined.
	unknown relationship = iota
	subtype
	supertype
)

const (
	startGlyph  = `╶──`
	branchGlyph = ` ├─`
	endGlyph    = ` └─`
	addGlyph    = ` + `
	indentGlyph = ` │ `
	blankGlyph  = `   `
)

type Node[T any] interface {
	comparable
	AddInherits(parent T) T
	Inherits() collections.SortedSet[T]

	// Implements determines if this interface implements the other interface.
	Implements(other T) bool
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
	comp  comp.Comparer[T]
	log   *logger.Logger
	count int
}

func New[T Node[T]](comp comp.Comparer[T], log *logger.Logger) Inheritance[T] {
	return &inheritanceImp[T]{
		roots: sortedSet.New(comp),
		comp:  comp,
		log:   log,
	}
}

func (in *inheritanceImp[T]) Process(node T) {
	relations := map[T]relationship{node: unknown}

	// pre-process all previously inherited nodes.
	for p := range node.Inherits().Enumerate().Seq() {
		in.Process(p)
		relations[p] = supertype
	}

	in.log.Logf(startGlyph+`(%d) insert: %v`, in.count, node)
	log2 := in.log.IndentWith(blankGlyph)
	in.addParent(in.roots, node, relations, log2)

	in.log.Logf(startGlyph+`(%d) done: %v`, in.count, node)
	for i, count := 0, node.Inherits().Count(); i < count; i++ {
		p := node.Inherits().Get(i)
		if i < count-1 {
			log2.Logf(branchGlyph+`(%d) %v`, i, p)
		} else {
			log2.Logf(endGlyph+`(%d) %v`, i, p)
		}
	}

	in.count++
}

func (in *inheritanceImp[T]) getRelationship(n, a T, relations map[T]relationship) (relationship, bool) {
	rel, touched := relations[a]
	if !touched {
		switch {
		case a.Implements(n):
			rel = subtype
		case n.Implements(a):
			rel = supertype
		default:
			rel = unknown
		}
		relations[a] = rel
	}
	return rel, touched
}

func (in *inheritanceImp[T]) addParent(siblings collections.SortedSet[T], n T, relations map[T]relationship, log *logger.Logger) {
	log2 := log.IndentWith(indentGlyph)
	addedToSibling := false
	parentedSiblings := false
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)

		rel, touched := in.getRelationship(n, a, relations)
		switch rel {
		case subtype:
			// Yi <: X, meaning `n` is a parent (sub-type) of `a`,
			// so add `n` as a parent of `a` and don't add it here.
			//
			// For example: {A, B} is a parent of {A, B, C} but {A, B, C}
			// may already have the parent {A} in it, so we have to recursively
			// call addParent to re-parent {A} as a parent of {A, B}.
			log.Logf(branchGlyph+`(%d) parent of %v`, i, a)
			if touched {
				log2.Log(endGlyph + ` skip: already checked`)
			} else {
				in.addParent(a.Inherits(), n, relations, log2)
			}
			addedToSibling = true

		case supertype:
			// Yi :> X, meaning `n` is a child (super-type) of `a`,
			// so move `a` from this set and add `a` as a parent of `n`.
			// This means that `n` is a parent in this set since otherwise
			// `a` would have been a parent to another object in this set.
			// We can simply add `a` to `n` since `a` has already been
			// checked against the other parents, hence it was in this set.
			log.Logf(branchGlyph+`(%d) child of %v`, i, a)
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
			log.Logf(branchGlyph+`(%d) else %v`, i, a)
			in.seekInherits(a.Inherits(), n, relations, log2)
		}
	}

	switch {
	case parentedSiblings:
		log.Log(endGlyph + ` add: parented sibling`)
		siblings.Add(n)

	case addedToSibling:
		log.Log(endGlyph + ` no-op: added to sibling`)

	default:
		log.Log(endGlyph + ` add: default`)
		siblings.Add(n)
	}
}

func (in *inheritanceImp[T]) seekInherits(siblings collections.SortedSet[T], n T, relations map[T]relationship, log *logger.Logger) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)

		rel, _ := in.getRelationship(n, a, relations)
		if rel == supertype {
			log.Logf(addGlyph+`%v`, a)
			n.AddInherits(a)
			continue
		}

		in.seekInherits(a.Inherits(), n, relations, log)
	}
}
