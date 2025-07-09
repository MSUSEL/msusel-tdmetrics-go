package constructs

import (
	"go/types"
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

// Construct is part of the source code.
type Construct interface {
	comp.Comparable[Construct]
	stringer.Stringerable

	// String gets a human readable string for this debugging this construct.
	String() string

	// Kind gets a string unique to each construct type.
	Kind() kind.Kind

	// Index gets the index of the construct, zero if unset.
	// The index will be set to the top level factory sorted set.
	Index() int

	// Duplicate is a flag set when this construct is identical to another
	// existing construct (one of the identical ones will not be marked as
	// a duplicate). Duplicates happen when one of the duplicates had a
	// reference, originally making it different, but once the reference was
	// replaced, the constraints were identical.
	Duplicate() bool

	// SetIndex sets the index of construct and indicates if it was a duplicate.
	SetIndex(index int, duplicate bool)

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
	// If required is true, then it will panic if a reference is not replicable.
	RemoveTempReferences(required bool) bool
}

// TempDeclRefContainer is any construct that can contain a temporary method reference.
type TempDeclRefContainer interface {
	Construct

	// RemoveTempDeclRefs should replace any found method reference with the
	// method that was referenced. References will already be looked up.
	// If required is true, then it will panic if a reference is not replicable.
	RemoveTempDeclRefs(required bool) bool
}

// NestType is a type that can be nested inside another type.
type NestType interface {
	Construct
	Name() string
	FuncType() *types.Func
}

// Nestable is a construct that can be nested inside another construct.
type Nestable interface {
	Construct
	Nest() NestType
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

	// These are constructs that can be have other constructs nested inside them.
	_ NestType = Method(nil)
	_ NestType = MethodInst(nil)

	// These are constructs that can be nested inside other constructs.
	_ Nestable = Object(nil)
	_ Nestable = InterfaceDecl(nil)
	_ Nestable = TempDeclRef(nil)
	_ Nestable = TempReference(nil)
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

func JsonSet[T Construct, S ~[]T](ctx *jsonify.Context, values S) *jsonify.List {
	return jsonify.NewSortedSetNonZero(ctx, values, Comparer[T]())
}

func ResolvedTempReference(td TypeDesc, required bool) (TypeDesc, bool) {
	changed := false
	if utils.IsNil(td) {
		panic(terror.New(`Construct given to ResolvedTempDeclRef was nil`))
	}
	resolved := td
	for {
		if resolved.Kind() == kind.TempReference {
			tr := resolved.(TempReference)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				if !required {
					return tr, changed
				}
				panic(terror.New(`TempReference in ResolvedTempReference resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, td))
			}
			changed = true
			continue
		}
		if resolved.Kind() == kind.TempTypeParamRef {
			tr := resolved.(TempTypeParamRef)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				if !required {
					return tr, changed
				}
				panic(terror.New(`TempTypeParamRef in ResolvedTempReference resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, td))
			}
			changed = true
			continue
		}
		break
	}
	return resolved, changed
}

func ResolvedTempDeclRef(con Construct, required bool) (Construct, bool) {
	changed := false
	if utils.IsNil(con) {
		panic(terror.New(`Construct given to ResolvedTempDeclRef was nil`))
	}
	resolved := con
	for {
		if resolved.Kind() == kind.TempReference {
			tr := resolved.(TempReference)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				if !required {
					return tr, changed
				}
				panic(terror.New(`TempReference in ResolvedTempDeclRef resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, con))
			}
			changed = true
			continue
		}
		if resolved.Kind() == kind.TempDeclRef {
			tr := resolved.(TempDeclRef)
			resolved = tr.ResolvedType()
			if utils.IsNil(resolved) {
				if !required {
					return tr, changed
				}
				panic(terror.New(`TempDeclRef in ResolvedTempDeclRef resolved to nil`).
					With(`Ref`, tr).
					With(`Start`, con))
			}
			changed = true
			continue
		}
		break
	}
	return resolved, changed
}

func ResolveTempDeclRefSet(set collections.SortedSet[Construct], required bool) bool {
	changed := false
	slice := slices.Clone(set.ToSlice())
	for i, s := range slice {
		ref, subChanged := ResolvedTempDeclRef(s, required)
		changed = changed || subChanged
		slice[i] = ref
	}
	assert.ArgHasNoNils(`resolved refs`, slice)
	if changed {
		set.Clear()
		set.Add(slice...)
	}
	return changed
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

func Cast[TOut, TIn Construct, S ~[]TIn](s S) []TOut {
	tps := make([]TOut, len(s))
	for i, tp := range s {
		tps[i] = any(tp).(TOut)
	}
	return tps
}

// ConstructCore is a shared data and methods for all constructs that
// may be embedded into a construct to quickly implement this data.
type ConstructCore struct {
	index     int
	duplicate bool
	alive     bool
}

func (c *ConstructCore) Index() int      { return c.index }
func (c *ConstructCore) Duplicate() bool { return c.duplicate }
func (c *ConstructCore) Alive() bool     { return c.alive }

func (c *ConstructCore) SetIndex(index int, duplicate bool) {
	c.index = index
	c.duplicate = duplicate
}

func (c *ConstructCore) SetAlive(alive bool) { c.alive = alive }
