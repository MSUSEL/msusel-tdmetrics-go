package constructs

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"maps"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface interface {
	TypeDesc
	_interface()

	IsSupertypeOf(other Interface) bool
	AddInheritors(inter Interface) bool
	SetInheritance()
}

type InterfaceArgs struct {
	RealType     *types.Interface
	Union        Union
	Methods      map[string]TypeDesc
	TypeParams   []Named
	InitInherits []Interface

	// Package is only needed if the real type is nil
	// so that a Go interface type has to be created.
	Package *packages.Package
}

func NewInterface(reg Register, args InterfaceArgs) Interface {
	methods := maps.Clone(args.Methods)

	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(errors.New(`must provide a package if the real type for an interface is nil`))
		}

		mTyp := []*types.Func{}
		pkg := args.Package.Types
		names := utils.SortedKeys(methods)
		for _, name := range names {
			sig := methods[name].GoType().(*types.Signature)
			if utils.IsNil(sig) {
				panic(fmt.Errorf(`nil signature for %s`, name))
			}
			f := types.NewFunc(token.NoPos, pkg, name, sig)
			mTyp = append(mTyp, f)
		}

		tEmb := []types.Type{}
		// TODO: Fill parameter types for interface.

		realType := types.NewInterfaceType(mTyp, tEmb)
		if realType == nil {
			panic(fmt.Errorf(`failed to create an interface`))
		}
		args.RealType = realType
	}

	return reg.RegisterInterface(&interfaceImp{
		realType:   args.RealType,
		typeParams: args.TypeParams,
		methods:    methods,
		union:      args.Union,
		inherits:   args.InitInherits,
	})
}

type interfaceImp struct {
	realType *types.Interface

	typeParams []Named
	methods    map[string]TypeDesc
	union      Union

	index      int
	inherits   []Interface
	inheritors []Interface
}

func (ti *interfaceImp) _interface() {}

func (ti *interfaceImp) Visit(v Visitor) {
	visitList(v, ti.typeParams)
	visitMap(v, ti.methods)
	visitTest(v, ti.union)
	visitList(v, ti.inherits)
	visitList(v, ti.inheritors)
}

func (ti *interfaceImp) SetIndex(index int) {
	ti.index = index
}

func (ti *interfaceImp) GoType() types.Type {
	return ti.realType
}

func (ti *interfaceImp) Equal(other TypeDesc) bool {
	return equalTest(ti, other, func(a, b *interfaceImp) bool {
		return equal(a.union, b.union) &&
			equalList(a.typeParams, b.typeParams) &&
			equalMap(a.methods, b.methods)
	})
}

func (ti *interfaceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ti.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `interface`).
		AddNonZero(ctx2, `typeParams`, ti.typeParams).
		AddNonZero(ctx2, `inherits`, ti.inherits).
		AddNonZero(ctx2, `union`, ti.union).
		AddNonZero(ctx2, `methods`, ti.methods)
}

func (ti *interfaceImp) String() string {
	return jsonify.ToString(ti)
}

func (ti *interfaceImp) IsSupertypeOf(other Interface) bool {
	otherTyp, ok := other.GoType().(*types.Interface)
	if !ok || utils.IsNil(ti.realType) || utils.IsNil(otherTyp) {
		// Baked in types don't have underlying interfaces
		// but also shouldn't be needed for any inheritance.
		return false
	}
	return types.Implements(ti.realType, otherTyp)
}

func (ti *interfaceImp) AddInheritors(other Interface) bool {
	inter, ok := other.(*interfaceImp)
	if !ok {
		return false
	}
	if ti == inter {
		return true
	}
	if !inter.IsSupertypeOf(ti) {
		return false
	}

	homed := false
	for _, other := range ti.inheritors {
		if other.AddInheritors(inter) {
			homed = true
		}
	}
	if homed {
		return true
	}

	changed := false
	for i, other := range ti.inheritors {
		if other.IsSupertypeOf(inter) {
			inter.inheritors = append(inter.inheritors, other)
			ti.inheritors[i] = nil
			changed = true
		}
	}
	if changed {
		ti.inheritors = utils.RemoveZeros(ti.inheritors)
	}

	ti.inheritors = append(ti.inheritors, inter)
	return true
}

func (ti *interfaceImp) SetInheritance() {
	for _, other := range ti.inheritors {
		otherInter, ok := other.(*interfaceImp)
		if !ok {
			panic(fmt.Sprintf(`interfaces in inheritors must be interfaceImps but got (%[1]T) %[1]v`, other))
		}
		otherInter.inherits = append(otherInter.inherits, ti)
	}
}
