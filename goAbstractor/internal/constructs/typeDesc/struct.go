package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Struct struct {
	Index  int
	Fields []*Field
}

func (ts *Struct) _isTypeDesc() {}

func (ts *Struct) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.GetBool(`onlyIndex`) {
		return jsonify.New(ctx, ts.Index)
	}

	ctx = ctx.Copy().Set(`onlyIndex`, true)

	return jsonify.NewMap().
		Add(ctx, `kind`, `struct`).
		Add(ctx, `fields`, ts.Fields)
}

type Field struct {
	Anonymous bool
	Name      string
	Type      TypeDesc
}

func (f *Field) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `anonymous`, f.Anonymous).
		AddNonZero(ctx, `name`, f.Name).
		Add(ctx, `type`, f.Type)
}
