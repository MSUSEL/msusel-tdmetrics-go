package typeDesc

import (
	"fmt"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type Interface struct {
	typ *types.Interface

	TypeParams []*Named
	Methods    map[string]TypeDesc
	Union      Union

	index      int
	Inherits   []*Interface
	Inheritors []*Interface
}

func NewInterface(typ *types.Interface) *Interface {
	return &Interface{
		typ:     typ,
		Methods: map[string]TypeDesc{},
	}
}

func (ti *Interface) SetIndex(index int) {
	ti.index = index
}

func (ti *Interface) GoType() types.Type {
	return ti.typ
}

func (ti *Interface) Equal(other TypeDesc) bool {
	return equalTest(ti, other, func(a, b *Interface) bool {
		return equal(a.Union, b.Union) &&
			equalList(a.TypeParams, b.TypeParams) &&
			equalMap(a.Methods, b.Methods)
	})
}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ti.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `interface`).
		AddNonZero(ctx2.Long(), `typeParams`, ti.TypeParams).
		AddNonZero(ctx2, `inherits`, ti.Inherits).
		AddNonZero(ctx2, `union`, ti.Union).
		AddNonZero(ctx2, `methods`, ti.Methods)
}

func (ti *Interface) String() string {
	return jsonify.ToString(ti)
}

func (ti *Interface) AddFunc(name string, sig TypeDesc) bool {
	if other, has := ti.Methods[name]; has {
		if other != sig {
			panic(fmt.Errorf(`function %v already exists with a different signature`, name))
		}
		return false
	}
	ti.Methods[name] = sig
	return true
}

func (ti *Interface) AddTypeParam(name string, t TypeDesc) *Named {
	tn := NewNamed(name, t)
	ti.TypeParams = append(ti.TypeParams, tn)
	return tn
}

func (ti *Interface) IsSupertypeOf(other *Interface) bool {
	if utils.IsNil(ti.typ) || utils.IsNil(other.typ) {
		// Baked in types don't have underlying interfaces
		// but also shouldn't be needed for any inheritance.
		return false
	}
	return types.Implements(ti.typ, other.typ)
}
