package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Struct struct {
	Fields     []*Named
	TypeParams []*Named

	Index int
	// Anonymous is the subset of fields that are anonymous.
	Anonymous []*Named
}

func NewStruct() *Struct {
	return &Struct{}
}

func (ts *Struct) _isTypeDesc() {}

func (ts *Struct) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.Index)
	}

	ctx2 := ctx.ShowKind().Long()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `struct`).
		Add(ctx2, `fields`, ts.Fields).
		AddNonZero(ctx2, `typeParams`, ts.TypeParams)
}

func (ts *Struct) String() string {
	return jsonify.ToString(ts)
}
