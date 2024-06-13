package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Struct interface {
	TypeDesc
	_struct()
}

type StructArgs struct {
	RealType *types.Struct
	Fields   []Named
}

func NewStruct(reg Register, args StructArgs) Struct {
	return reg.RegisterStruct(&structImp{
		realType: args.RealType,
		fields:   args.Fields,
	})
}

type structImp struct {
	realType *types.Struct
	fields   []Named
	index    int
}

func (ts *structImp) _struct() {}

func (ts *structImp) Visit(v Visitor) {
	visitList(v, ts.fields)
}

func (ts *structImp) SetIndex(index int) {
	ts.index = index
}

func (ts *structImp) GoType() types.Type {
	return ts.realType
}

func (ts *structImp) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *structImp) bool {
		return equalList(a.fields, b.fields)
	})
}

func (ts *structImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `struct`).
		AddNonZero(ctx2, `fields`, ts.fields)
}

func (ts *structImp) String() string {
	return jsonify.ToString(ts)
}
