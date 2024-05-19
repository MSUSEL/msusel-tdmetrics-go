package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Solid represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Solid struct {
	typ types.Type

	Target     TypeDesc
	TypeParams []TypeDesc
}

func NewSolid(typ types.Type, target TypeDesc, tp ...TypeDesc) *Solid {
	return &Solid{
		typ:        typ,
		Target:     target,
		TypeParams: tp,
	}
}

func (ts *Solid) GoType() types.Type {
	return ts.typ
}

func (ts *Solid) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *Solid) bool {
		return equal(a.Target, b.Target) &&
			equalList(a.TypeParams, b.TypeParams)
	})
}

func (ts *Solid) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `solid`).
		AddNonZero(ctx.ShowKind().Short(), `target`, ts.Target).
		AddNonZero(ctx.HideKind().Short(), `typeParams`, ts.TypeParams)
}

func (ts *Solid) String() string {
	return jsonify.ToString(ts)
}

func (ts *Solid) AppendTypeParam(tp ...TypeDesc) {
	ts.TypeParams = append(ts.TypeParams, tp...)
}
