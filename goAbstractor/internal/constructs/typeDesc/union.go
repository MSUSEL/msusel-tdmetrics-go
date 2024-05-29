package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// TODO: should this be expanded so that it can be removed?
// e.g. `Foo[T string|int]` becomes `Foo_string` and `Foo_int`.
//
// Exact types are like `string|int|bool` where the type must match exactly.
// Approx types are like `~string|~int` where they type may be exact or
// an extension of the base type.
type Union interface {
	TypeDesc

	AddType(approx bool, td ...TypeDesc)
}

func NewUnion(typ *types.Union) Union {
	return &unionImp{
		typ: typ,
	}
}

type unionImp struct {
	typ *types.Union

	index  int
	exact  []TypeDesc
	approx []TypeDesc
}

func (t *unionImp) SetIndex(index int) {
	t.index = index
}

func (t *unionImp) GoType() types.Type {
	return t.typ
}

func (t *unionImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *unionImp) bool {
		return equalList(a.exact, b.exact) &&
			equalList(a.approx, b.approx)
	})
}

func (t *unionImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `union`).
		AddNonZero(ctx.ShowKind().Short(), `exact`, t.exact).
		AddNonZero(ctx.ShowKind().Short(), `approx`, t.approx)
}

func (t *unionImp) String() string {
	return jsonify.ToString(t)
}

func (t *unionImp) AddType(approx bool, td ...TypeDesc) {
	if approx {
		t.approx = append(t.approx, td...)
	} else {
		t.exact = append(t.exact, td...)
	}
}
