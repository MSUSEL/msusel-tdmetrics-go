package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

// TODO: should this be expanded so that it can be removed?
// e.g. `Foo[T string|int]` becomes `Foo_string` and `Foo_int`.
type Union struct {
	Types []TypeDesc
}

func NewUnion(types ...TypeDesc) *Union {
	return &Union{Types: types}
}

func (t *Union) _isTypeDesc() {}

func (t *Union) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `union`).
		AddNonZero(ctx.ShowKind().Short(), `types`, t.Types)
}

func (t *Union) String() string {
	return jsonify.ToString(t)
}

func (t *Union) AppendType(td ...TypeDesc) {
	t.Types = append(t.Types, td...)
}
