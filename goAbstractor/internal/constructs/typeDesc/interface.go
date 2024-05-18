package typeDesc

import (
	"fmt"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Interface struct {
	typ *types.Interface

	TypeParams []*Named
	Methods    map[string]TypeDesc
	Union      *Union

	Index      int
	Inherits   []*Interface
	Inheritors []*Interface
}

func NewInterface(typ *types.Interface) *Interface {
	return &Interface{
		typ:     typ,
		Methods: map[string]TypeDesc{},
	}
}

func (ti *Interface) GoType() types.Type {
	return ti.typ
}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ti.Index)
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

func (ti *Interface) HasFunc(name string, sig TypeDesc) bool {
	other, has := ti.Methods[name]
	// TODO: Need to handle types which have been made solid?
	// e.g. `Foo[T](val T)` with `T` as `int` and `Bar(val int)`
	//
	// See https://go.dev/doc/tutorial/generics
	// See https://stackoverflow.com/questions/70888240/whats-the-meaning-of-the-new-tilde-token-in-go
	//
	// DON'T DO ON YOUR OWN, Use reflection and the Implements methods.
	return has && sig == other
}

func (ti *Interface) IsSupertypeOf(other *Interface) bool {
	for name, sig := range other.Methods {
		if !ti.HasFunc(name, sig) {
			return false
		}
	}
	return true
}
