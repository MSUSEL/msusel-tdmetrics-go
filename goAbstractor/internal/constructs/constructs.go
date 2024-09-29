package constructs

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

// Construct is part of the source code.
type Construct interface {
	comp.Comparable[Construct]

	// String gets a human readable string for this debugging this construct.
	String() string

	// Kind gets a string unique to each construct type.
	Kind() kind.Kind

	// Index gets the index of the construct, zero if unset.
	// The index will be set to the top level factory sorted set.
	Index() int

	// SetIndex sets the index of construct.
	SetIndex(index int)

	// Alive indicates that this construct is reachable
	// from any entry point in the compiled project.
	Alive() bool

	// SetAlive sets if the given construct is alive.
	SetAlive(alive bool)
}

// TempReferenceContainer is any construct that can contain a temporary reference.
type TempReferenceContainer interface {
	Construct

	// RemoveTempReferences should replace any found reference with the type
	// description that was referenced. References will already be looked up.
	RemoveTempReferences()
}

// TempDeclRefContainer is any construct that can contain a temporary method reference.
type TempDeclRefContainer interface {
	Construct

	// RemoveTempDeclRefs should replace any found method reference with the
	// method that was referenced. References will already be looked up.
	RemoveTempDeclRefs()
}

var (
	_ Construct = Abstract(nil)
	_ Construct = Argument(nil)
	_ Construct = Field(nil)
	_ Construct = Package(nil)
	_ Construct = Metrics(nil)
	_ Construct = MethodInst(nil)

	// These are the implementations of type descriptions.
	// None of these have generics defined on them but may carry
	// type parameters for the generic declaration that they are part of.
	_ TypeDesc = Basic(nil)
	_ TypeDesc = InterfaceDesc(nil)
	_ TypeDesc = InterfaceInst(nil)
	_ TypeDesc = ObjectInst(nil)
	_ TypeDesc = Signature(nil)
	_ TypeDesc = StructDesc(nil)
	_ TypeDesc = TypeParam(nil)
	_ TypeDesc = TempReference(nil)

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
	if c := strings.Compare(string(a.Kind()), string(b.Kind())); c != 0 {
		return c
	}
	return cmp(a, b.(T))
}

func ResolvedTempReference(td TypeDesc) TypeDesc {
	for td.Kind() == kind.TempReference {
		td = td.(TempReference).ResolvedType()
	}
	return td
}

func ResolvedTempDeclRef(con Construct) Construct {
	for {
		if con.Kind() == kind.TempReference {
			con = con.(TempReference).ResolvedType()
			continue
		}
		if con.Kind() == kind.TempDeclRef {
			con = con.(TempDeclRef).ResolvedType()
			continue
		}
		break
	}
	return con
}

func BlankName(name string) bool {
	return len(name) <= 0 || name == `_` || name == `.`
}
