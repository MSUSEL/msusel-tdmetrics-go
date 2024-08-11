package structDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `structDesc`

type StructDesc interface {
	typeDescs.TypeDesc
	_structDesc()
}

type Args struct {
	RealType types.Type

	Fields []field.Field
}

type structDescImp struct {
	realType types.Type

	fields []field.Field

	index int
}

func newStructDesc(args Args) StructDesc {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNoNils(`fields`, args.Fields)

	return &structDescImp{
		realType: args.RealType,
		fields:   args.Fields,
	}
}

func (d *structDescImp) _structDesc()       {}
func (d *structDescImp) Kind() string       { return Kind }
func (d *structDescImp) SetIndex(index int) { d.index = index }
func (d *structDescImp) GoType() types.Type { return d.realType }

func (d *structDescImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[StructDesc](d, other, Comparer())
}

func Comparer() comp.Comparer[StructDesc] {
	return func(a, b StructDesc) int {
		aImp, bImp := a.(*structDescImp), b.(*structDescImp)
		return constructs.SliceComparer[field.Field]()(aImp.fields, bImp.fields)
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
