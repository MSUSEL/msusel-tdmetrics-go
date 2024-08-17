package structDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type structDescImp struct {
	realType types.Type

	fields []constructs.Field

	index int
}

func newStructDesc(args constructs.StructDescArgs) constructs.StructDesc {
	assert.ArgHasNoNils(`fields`, args.Fields)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		// TODO: Implement
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &structDescImp{
		realType: args.RealType,
		fields:   args.Fields,
	}
}

func (d *structDescImp) IsTypeDesc()        {}
func (d *structDescImp) IsStructDesc()      {}
func (d *structDescImp) Kind() kind.Kind    { return kind.StructDesc }
func (d *structDescImp) SetIndex(index int) { d.index = index }
func (d *structDescImp) GoType() types.Type { return d.realType }

func (d *structDescImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.StructDesc](d, other, Comparer())
}

func Comparer() comp.Comparer[constructs.StructDesc] {
	return func(a, b constructs.StructDesc) int {
		aImp, bImp := a.(*structDescImp), b.(*structDescImp)
		return constructs.SliceComparer[constructs.Field]()(aImp.fields, bImp.fields)
	}
}

func (d *structDescImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, d.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, d.index).
		AddNonZero(ctx2, `fields`, d.fields)
}
