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

	ctx2 := ctx.Copy().
		Remove(`noKind`).
		Set(`onlyIndex`, true)

	showKind := !ctx.GetBool(`noKind`)
	return jsonify.NewMap().
		AddIf(ctx2, showKind, `kind`, `struct`).
		Add(ctx2, `fields`, ts.Fields)
}

func (ts *Struct) String() string {
	return jsonify.ToString(ts)
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

func (f *Field) String() string {
	return jsonify.ToString(f)
}
