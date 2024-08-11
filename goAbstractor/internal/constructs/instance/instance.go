package instance

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type instanceImp struct {
	resolved   constructs.TypeDesc
	typeParams []constructs.TypeDesc

	index int
}

func newInstance(args constructs.InstanceArgs) constructs.Instance {
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`type params`, args.TypeParams)
	assert.ArgNoNils(`type params`, args.TypeParams)
	return &instanceImp{
		resolved:   args.Resolved,
		typeParams: args.TypeParams,
	}
}

func (i *instanceImp) IsInstance()        {}
func (i *instanceImp) Kind() kind.Kind    { return kind.Instance }
func (i *instanceImp) SetIndex(index int) { i.index = index }

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Instance](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Instance] {
	return func(a, b constructs.Instance) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		return comp.Or(
			constructs.ComparerPend(aImp.resolved, bImp.resolved),
			constructs.SliceComparerPend(bImp.typeParams, bImp.typeParams),
		)
	}
}

func (i *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, i.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, i.index).
		AddNonZero(ctx2, `resolved`, i.resolved).
		AddNonZero(ctx2, `typeParams`, i.typeParams)
}
