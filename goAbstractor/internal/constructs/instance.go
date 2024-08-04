package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Instance represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Instance interface {
	Construct
	_instance()
}

type InstanceArgs struct {
	RealType   types.Type
	Resolved   TypeDesc
	TypeParams []TypeDesc
}

type instanceImp struct {
	realType   types.Type
	resolved   TypeDesc
	typeParams []TypeDesc

	index int
}

func newInstance(args InstanceArgs) Instance {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`type params`, args.TypeParams)
	assert.ArgNoNils(`type params`, args.TypeParams)

	return &instanceImp{
		realType:   args.RealType,
		resolved:   args.Resolved,
		typeParams: args.TypeParams,
	}
}

func (i *instanceImp) _instance()         {}
func (i *instanceImp) Kind() kind.Kind    { return kind.Instance }
func (i *instanceImp) setIndex(index int) { i.index = index }

func (i *instanceImp) compareTo(other Construct) int {
	b := other.(*instanceImp)
	return or(
		func() int { return Compare(i.resolved, b.resolved) },
		func() int { return compareSlice(i.typeParams, b.typeParams) },
	)
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
