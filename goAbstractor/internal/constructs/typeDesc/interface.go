package typeDesc

import (
	"go/types"
	"maps"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface interface {
	TypeDesc

	IsSupertypeOf(other Interface) bool
	AppendInherits(inherits ...Interface)
	AddInheritors(inter Interface) bool
	SetInheritance()
}

func NewInterface(typ *types.Interface, union Union, methods map[string]TypeDesc, typeParams ...Named) Interface {
	return &interfaceImp{
		typ:        typ,
		typeParams: typeParams,
		methods:    maps.Clone(methods),
		union:      union,
	}
}

type interfaceImp struct {
	typ *types.Interface

	typeParams []Named
	methods    map[string]TypeDesc
	union      Union

	index      int
	inherits   []Interface
	inheritors []Interface
}

func (ti *interfaceImp) SetIndex(index int) {
	ti.index = index
}

func (ti *interfaceImp) GoType() types.Type {
	return ti.typ
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
		AddNonZero(ctx2.Long(), `typeParams`, ti.typeParams).
		AddNonZero(ctx2, `inherits`, ti.inherits).
		AddNonZero(ctx2, `union`, ti.union).
		AddNonZero(ctx2, `methods`, ti.methods)
}

func (ti *interfaceImp) String() string {
	return jsonify.ToString(ti)
}

func (ti *interfaceImp) IsSupertypeOf(other Interface) bool {
	otherTyp, ok := other.GoType().(*types.Interface)
	if !ok || utils.IsNil(ti.typ) || utils.IsNil(otherTyp) {
		// Baked in types don't have underlying interfaces
		// but also shouldn't be needed for any inheritance.
		return false
	}
	return types.Implements(ti.typ, otherTyp)
}

func (ti *interfaceImp) AppendInherits(inherits ...Interface) {
	ti.inherits = append(ti.inherits, inherits...)
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
		other.AppendInherits(ti)
	}
}
