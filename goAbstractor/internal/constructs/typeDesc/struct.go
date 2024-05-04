package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Struct struct {
	Fields     []*Named
	TypeParams []*Named

	Index int
	// Embedded is the subset of fields that are embedded.
	Embedded []*Named
}

func NewStruct() *Struct {
	return &Struct{}
}

func (ts *Struct) _isTypeDesc() {}

func (ts *Struct) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.Index)
	}

	ctx2 := ctx.HideKind().Long()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `struct`).
		Add(ctx2, `fields`, ts.Fields).
		AddNonZero(ctx2, `typeParams`, ts.TypeParams)
}

func (ts *Struct) String() string {
	return jsonify.ToString(ts)
}

func (ts *Struct) AddField(name string, t TypeDesc, embedded bool) *Named {
	tn := NewNamed(name, t)
	ts.Fields = append(ts.Fields, tn)
	if embedded {
		ts.Embedded = append(ts.Embedded, tn)
	}
	return tn
}

func (ts *Struct) AddTypeParam(name string, t *Interface) *Named {
	tn := NewNamed(name, t)
	ts.TypeParams = append(ts.TypeParams, tn)
	return tn
}

func (ts *Struct) AppendTypeParam(tp ...*Named) {
	ts.TypeParams = append(ts.TypeParams, tp...)
}
