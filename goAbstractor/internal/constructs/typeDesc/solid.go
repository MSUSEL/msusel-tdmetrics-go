package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

// Solid represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Solid struct {
	Target     TypeDesc
	TypeParams []TypeDesc
}

func NewSolid(target TypeDesc, tp ...TypeDesc) *Solid {
	return &Solid{
		Target:     target,
		TypeParams: tp,
	}
}

func (ts *Solid) _isTypeDesc() {}

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
