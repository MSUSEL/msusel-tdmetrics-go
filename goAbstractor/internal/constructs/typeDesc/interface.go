package typeDesc

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface interface {
	TypeDesc

	SetUnion(union Union)
	AddFunc(name string, sig TypeDesc) bool
	AddTypeParam(name string, t TypeDesc) Named
	IsSupertypeOf(other Interface) bool
	AppendInherits(inherits ...Interface)
	AddInheritors(inter Interface) bool
	SetInheritance()
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

func NewInterface(typ *types.Interface) Interface {
	return &interfaceImp{
		typ:     typ,
		methods: map[string]TypeDesc{},
	}
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

func (ti *interfaceImp) SetUnion(union Union) {
	ti.union = union
}

func (ti *interfaceImp) AddFunc(name string, sig TypeDesc) bool {
	if other, has := ti.methods[name]; has {
		if other != sig {
			panic(fmt.Errorf(`function %v already exists with a different signature`, name))
		}
		return false
	}
	ti.methods[name] = sig
	return true
}

func (ti *interfaceImp) AddTypeParam(name string, t TypeDesc) Named {
	tn := NewNamed(name, t)
	ti.typeParams = append(ti.typeParams, tn)
	return tn
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
