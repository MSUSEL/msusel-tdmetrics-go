package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

// Construct is part of the source code.
type Construct interface {
	comp.Comparable[Construct]

	// Kind gets a string unique to each construct type.
	Kind() kind.Kind

	// setIndex sets the unique index of construct.
	// Indices will be 1 based so that 0 is unset.
	SetIndex(index int)
}

var (
	_ Construct = Abstract(nil)
	_ Construct = Argument(nil)
	_ Construct = Field(nil)
	_ Construct = Package(nil)

	// These are the implementations of type descriptions.
	// None of these have generics defined on them but may carry
	// type parameters for the generic declaration that they are part of.
	_ TypeDesc = Basic(nil)
	_ TypeDesc = Instance(nil)
	_ TypeDesc = InterfaceDesc(nil)
	_ TypeDesc = Reference(nil)
	_ TypeDesc = Signature(nil)
	_ TypeDesc = StructDesc(nil)
	_ TypeDesc = TypeParam(nil)

	// A TypeDecl is both a TypeDesc and a Declaration.
	_ TypeDesc    = TypeDecl(nil)
	_ Declaration = TypeDecl(nil)

	// These are TypeDecls. They are declarations and also type descriptions
	// because they can be used by name, i.e. `var X ObjectFoo`.
	_ TypeDecl = Object(nil)
	_ TypeDecl = InterfaceDecl(nil)

	// These are type declarations only. They can not be used at TypeDesc.
	_ Declaration = Method(nil)
	_ Declaration = Value(nil)
)

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
	if utils.IsNil(a) {
		return utils.Ternary(utils.IsNil(b), 0, -1)
	}
	if utils.IsNil(b) {
		return 1
	}
	return comp.Or(
		comp.DefaultPend(string(a.Kind()), string(b.Kind())),
		cmp.Pend(a, b.(T)),
	)
}
