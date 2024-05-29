package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Struct interface {
	TypeDesc

	AddField(name string, t TypeDesc, embedded bool) Named
	AppendField(embedded bool, fields ...Named)
	AddTypeParam(name string, t Interface) Named
	AppendTypeParam(tp ...Named)
}

func NewStruct(typ *types.Struct) Struct {
	return &structImp{
		typ: typ,
	}
}

type structImp struct {
	typ *types.Struct

	fields     []Named
	typeParams []Named

	index int

	// embedded is the subset of fields that are embedded.
	embedded []Named
}

func (ts *structImp) SetIndex(index int) {
	ts.index = index
}

func (ts *structImp) GoType() types.Type {
	return ts.typ
}

func (ts *structImp) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *structImp) bool {
		return equalList(a.fields, b.fields) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (ts *structImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.index)
	}

	ctx2 := ctx.HideKind().Long()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `struct`).
		Add(ctx2, `fields`, ts.fields).
		AddNonZero(ctx2, `typeParams`, ts.typeParams)
}

func (ts *structImp) String() string {
	return jsonify.ToString(ts)
}

func (ts *structImp) AddField(name string, t TypeDesc, embedded bool) Named {
	tn := NewNamed(name, t)
	ts.AppendField(embedded, tn)
	return tn
}

func (ts *structImp) AppendField(embedded bool, fields ...Named) {
	ts.fields = append(ts.fields, fields...)
	if embedded {
		ts.embedded = append(ts.embedded, fields...)
	}
}

func (ts *structImp) AddTypeParam(name string, t Interface) Named {
	tn := NewNamed(name, t)
	ts.typeParams = append(ts.typeParams, tn)
	return tn
}

func (ts *structImp) AppendTypeParam(tp ...Named) {
	ts.typeParams = append(ts.typeParams, tp...)
}
