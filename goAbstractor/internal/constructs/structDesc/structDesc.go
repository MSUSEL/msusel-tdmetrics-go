package structDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `structDesc`

type Args struct {
	RealType types.Type

	Fields []components.Field
}

type structDescImp struct {
	realType types.Type

	fields []components.Field

	index int
}

func New(args Args) typeDescs.StructDesc {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNoNils(`fields`, args.Fields)

	return &structDescImp{
		realType: args.RealType,
		fields:   args.Fields,
	}
}

func (d *structDescImp) IsTypeDesc()        {}
func (d *structDescImp) IsStructDesc()      {}
func (d *structDescImp) Kind() string       { return Kind }
func (d *structDescImp) SetIndex(index int) { d.index = index }
func (d *structDescImp) GoType() types.Type { return d.realType }

func (d *structDescImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[typeDescs.StructDesc](d, other, Comparer())
}

func Comparer() comp.Comparer[typeDescs.StructDesc] {
	return func(a, b typeDescs.StructDesc) int {
		aImp, bImp := a.(*structDescImp), b.(*structDescImp)
		return constructs.SliceComparer[components.Field]()(aImp.fields, bImp.fields)
	}
}

func (d *structDescImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, d.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, d.index).
		AddNonZero(ctx2, `fields`, d.fields)
}
