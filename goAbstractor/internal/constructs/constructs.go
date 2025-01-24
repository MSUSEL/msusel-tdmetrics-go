package constructs

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
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
	// This will also remove any TempTypeParamRefs.
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
	_ Construct = TempDeclRef(nil)

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
	_ TypeDesc = TempTypeParamRef(nil)

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
	if utils.IsNil(td) {
		panic(terror.New(`Construct given to ResolvedTempDeclRef was nil`))
	}
	resolved := td
	for {
		if resolved.Kind() == kind.TempReference {
			tr := resolved.(TempReference)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				panic(terror.New(`TempReference in ResolvedTempReference resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, td))
			}
			continue
		}
		if resolved.Kind() == kind.TempTypeParamRef {
			tr := resolved.(TempTypeParamRef)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				panic(terror.New(`TempTypeParamRef in ResolvedTempReference resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, td))
			}
			continue
		}
		break
	}
	return resolved
}

func ResolvedTempDeclRef(con Construct) Construct {
	if utils.IsNil(con) {
		panic(terror.New(`Construct given to ResolvedTempDeclRef was nil`))
	}
	resolved := con
	for {
		if resolved.Kind() == kind.TempReference {
			tr := resolved.(TempReference)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				panic(terror.New(`TempReference in ResolvedTempDeclRef resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, con))
			}
			continue
		}
		if resolved.Kind() == kind.TempDeclRef {
			tr := resolved.(TempDeclRef)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				panic(terror.New(`TempDeclRef in ResolvedTempDeclRef resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, con))
			}
			continue
		}
		break
	}
	return resolved
}

func BlankName(name string) bool {
	return len(name) <= 0 || name == `_` || name == `.`
}

func FindSigByName(abs []Abstract, name string) Signature {
	for _, ab := range abs {
		if ab.Name() == name {
			return ab.Signature()
		}
	}
	panic(terror.New(`failed to find signature in interface abstract by name`).
		With(`name`, name).
		With(`abs`, abs))
}
