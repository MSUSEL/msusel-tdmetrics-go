package constructs

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Interface interface {
		TypeDesc
		_interface()

		IsSupertypeOf(other Interface) bool
		AddInheritors(inter Interface) bool
		SetInheritance()
	}

	InterfaceArgs struct {
		RealType   *types.Interface
		Union      Union
		Methods    []Named
		TypeParams []Named

		// Package is only needed if the real type is nil
		// so that a Go interface type has to be created.
		Package *packages.Package
	}

	interfaceImp struct {
		realType *types.Interface

		typeParams []Named
		methods    []Named
		union      Union

		index      int
		inherits   []Interface
		inheritors []Interface
	}
)

func newInterface(args InterfaceArgs) Interface {
	methods := slices.Clone(args.Methods)
	tp := slices.Clone(args.TypeParams)

	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(errors.New(`must provide a package if the real type for an interface is nil`))
		}

		mTyp := []*types.Func{}
		pkg := args.Package.Types
		for _, named := range methods {
			name := named.Name()
			sig := named.Type().GoType().(*types.Signature)
			if utils.IsNil(sig) {
				panic(fmt.Errorf(`nil signature for %s`, name))
			}
			f := types.NewFunc(token.NoPos, pkg, name, sig)
			mTyp = append(mTyp, f)
		}

		tEmb := make([]types.Type, len(args.TypeParams))
		for i, n := range args.TypeParams {
			tEmb[i] = n.GoType()
		}
		if !utils.IsNil(args.Union) {
			tEmb = append(tEmb, args.Union.GoType())
		}

		realType := types.NewInterfaceType(mTyp, tEmb)
		if realType == nil {
			panic(fmt.Errorf(`failed to create an interface`))
		}
		args.RealType = realType
	}

	return &interfaceImp{
		realType:   args.RealType,
		typeParams: tp,
		methods:    methods,
		union:      args.Union,
	}
}

func (it *interfaceImp) _interface()        {}
func (it *interfaceImp) Kind() kind.Kind    { return kind.Interface }
func (it *interfaceImp) SetIndex(index int) { it.index = index }
func (it *interfaceImp) GoType() types.Type { return it.realType }

func (it *interfaceImp) Visit(v visitor.Visitor) bool {
	return visitor.Visit(v, it.typeParams...) &&
		visitor.Visit(v, it.methods...) &&
		visitor.Visit(v, it.union) &&
		visitor.Visit(v, it.inherits...)
}

func (it *interfaceImp) CompareTo(other Construct) int {
	b := other.(*interfaceImp)
	if cmp := Compare(it.union, b.union); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(it.typeParams, b.typeParams); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(it.methods, b.methods); cmp != 0 {
		return cmp
	}
	return 0
}

func (it *interfaceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, it.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, it.Kind()).
		AddNonZero(ctx2, `typeParams`, it.typeParams).
		AddNonZero(ctx2, `inherits`, it.inherits).
		AddNonZero(ctx2, `union`, it.union).
		AddNonZero(ctx2, `methods`, it.methods)
}

func (it *interfaceImp) IsSupertypeOf(other Interface) bool {
	otherIt, ok := other.GoType().(*types.Interface)
	if !ok || utils.IsNil(it.realType) || utils.IsNil(otherIt) {
		// Baked in types don't have underlying interfaces
		// but also shouldn't be needed for any inheritance.
		return false
	}
	return types.Implements(it.realType, otherIt)
}

func (it *interfaceImp) AddInheritors(other Interface) bool {
	inter, ok := other.(*interfaceImp)
	if !ok {
		return false
	}
	if it == inter {
		return true
	}
	if !inter.IsSupertypeOf(it) {
		return false
	}

	homed := false
	for _, other := range it.inheritors {
		if other.AddInheritors(inter) {
			homed = true
		}
	}
	if homed {
		return true
	}

	changed := false
	for i, other := range it.inheritors {
		if other.IsSupertypeOf(inter) {
			inter.inheritors = append(inter.inheritors, other)
			it.inheritors[i] = nil
			changed = true
		}
	}
	if changed {
		it.inheritors = utils.RemoveZeros(it.inheritors)
	}

	it.inheritors = append(it.inheritors, inter)
	return true
}

func (it *interfaceImp) SetInheritance() {
	for _, other := range it.inheritors {
		otherInter, ok := other.(*interfaceImp)
		if !ok {
			panic(fmt.Sprintf(`interfaces in inheritors must be interfaceImps but got (%[1]T) %[1]v`, other))
		}
		otherInter.inherits = append(otherInter.inherits, it)
	}
}
