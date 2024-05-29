package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Solid represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Solid interface {
	TypeDesc

	AppendTypeParam(tp ...TypeDesc)
}

func NewSolid(typ types.Type, target TypeDesc, tp ...TypeDesc) Solid {
	return &solidImp{
		typ:        typ,
		target:     target,
		typeParams: tp,
	}
}

type solidImp struct {
	typ types.Type

	index      int
	target     TypeDesc
	typeParams []TypeDesc
}

func (ts *solidImp) SetIndex(index int) {
	ts.index = index
}

func (ts *solidImp) GoType() types.Type {
	return ts.typ
}

func (ts *solidImp) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *solidImp) bool {
		return equal(a.target, b.target) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (ts *solidImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.index)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `solid`).
		AddNonZero(ctx.ShowKind().Short(), `target`, ts.target).
		AddNonZero(ctx.ShowKind().Short(), `typeParams`, ts.typeParams)
}

func (ts *solidImp) String() string {
	return jsonify.ToString(ts)
}

func (ts *solidImp) AppendTypeParam(tp ...TypeDesc) {
	ts.typeParams = append(ts.typeParams, tp...)
}
