package instance

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `instance`

// Instance represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Instance interface {
	constructs.Construct
	_instance()
}

type Args struct {
	Resolved   typeDescs.TypeDesc
	TypeParams []typeDescs.TypeDesc
}

type instanceImp struct {
	resolved   typeDescs.TypeDesc
	typeParams []typeDescs.TypeDesc

	index int
}

func newInstance(args Args) Instance {
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`type params`, args.TypeParams)
	assert.ArgNoNils(`type params`, args.TypeParams)
	return &instanceImp{
		resolved:   args.Resolved,
		typeParams: args.TypeParams,
	}
}

func (i *instanceImp) _instance()         {}
func (i *instanceImp) Kind() string       { return Kind }
func (i *instanceImp) SetIndex(index int) { i.index = index }

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[Instance](i, other, Comparer())
}

func Comparer() comp.Comparer[Instance] {
	return func(a, b Instance) int {
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
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, i.index).
		AddNonZero(ctx2, `resolved`, i.resolved).
		AddNonZero(ctx2, `typeParams`, i.typeParams)
}
