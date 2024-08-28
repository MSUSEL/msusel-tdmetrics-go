package structDesc

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
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

		fields := make([]*types.Var, len(args.Fields))
		for i, field := range args.Fields {
			fields[i] = types.NewField(token.NoPos, args.Package.Types,
				field.Name(), field.Type().GoType(), field.Embedded())
		}
		args.RealType = types.NewStruct(fields, nil)
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &structDescImp{
		realType: args.RealType,
		fields:   args.Fields,
	}
}

func (d *structDescImp) IsTypeDesc()   {}
func (d *structDescImp) IsStructDesc() {}

func (d *structDescImp) Kind() kind.Kind    { return kind.StructDesc }
func (d *structDescImp) Index() int         { return d.index }
func (d *structDescImp) SetIndex(index int) { d.index = index }
func (d *structDescImp) GoType() types.Type { return d.realType }

func (d *structDescImp) Fields() []constructs.Field { return d.fields }

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
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, d.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, d.Kind(), d.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, d.index).
		AddNonZero(ctx.OnlyIndex(), `fields`, d.fields)
}

func (d *structDescImp) String() string {
	return `struct{ ` + enumerator.Enumerate(d.fields...).Join(`; `) + `}`
}
