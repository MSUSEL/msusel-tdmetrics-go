package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Struct struct {
	Fields []*Field

	Index int
}

func NewStruct() *Struct {
	return &Struct{}
}

func (ts *Struct) _isTypeDesc() {}

func (ts *Struct) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.Short {
		return jsonify.New(ctx, ts.Index)
	}

	ctx2 := ctx.Copy()
	ctx2.NoKind = false
	ctx2.Short = true

	return jsonify.NewMap().
		AddIf(ctx2, !ctx.NoKind, `kind`, `struct`).
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
