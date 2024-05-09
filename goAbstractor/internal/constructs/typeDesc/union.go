package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

// TODO: should this be expanded so that it can be removed?
// e.g. `Foo[T string|int]` becomes `Foo_string` and `Foo_int`.
//
// Exact types are like `string|int|bool` where the type must match exactly.
// Approx types are like `~string|~int` where they type may be exact or
// an extension of the base type.
type Union struct {
	Exact  []TypeDesc
	Approx []TypeDesc
}

func NewUnion() *Union {
	return &Union{}
}

func (t *Union) _isTypeDesc() {}

func (t *Union) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `union`).
		AddNonZero(ctx.ShowKind().Short(), `exact`, t.Exact).
		AddNonZero(ctx.ShowKind().Short(), `approx`, t.Approx)
}

func (t *Union) String() string {
	return jsonify.ToString(t)
}

func (t *Union) AddType(approx bool, td ...TypeDesc) {
	if approx {
		t.Approx = append(t.Approx, td...)
	} else {
		t.Exact = append(t.Exact, td...)
	}
}
