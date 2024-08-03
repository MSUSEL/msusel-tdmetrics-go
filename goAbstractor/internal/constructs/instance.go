package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	// Instance represents a generic type that has been resolved to a specific type
	// with specific type parameters, e.g. List<T> might be resolved to List<int>.
	// The type parameter resolution may be referencing another type parameter,
	// e.g. a method signature inside a generic interface.
	Instance interface {
		TypeDesc
		_instance()
	}

	InstanceArgs struct {
		RealType   types.Type
		Target     TypeDesc
		TypeParams []TypeDesc
	}

	instanceImp struct {
		realType types.Type

		index      int
		target     TypeDesc
		typeParams []TypeDesc
	}
)

func newInstance(args InstanceArgs) Instance {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`target`, args.Target)
	assert.ArgNotEmpty(`type params`, args.TypeParams)
	assert.ArgNoNils(`type params`, args.TypeParams)

	return &instanceImp{
		realType:   args.RealType,
		target:     args.Target,
		typeParams: args.TypeParams,
	}
}

func (s *instanceImp) _instance()         {}
func (s *instanceImp) Kind() kind.Kind    { return kind.Solid }
func (s *instanceImp) SetIndex(index int) { s.index = index }
func (s *instanceImp) GoType() types.Type { return s.realType }

func (s *instanceImp) CompareTo(other Construct) int {
	b := other.(*instanceImp)
	if cmp := Compare(s.target, b.target); cmp != 0 {
		return cmp
	}
	return CompareSlice(s.typeParams, b.typeParams)
}

func (s *instanceImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, s.target)
	visitor.Visit(v, s.typeParams...)
}

func (s *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, s.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, s.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, s.index).
		AddNonZero(ctx2, `target`, s.target).
		AddNonZero(ctx2, `typeParams`, s.typeParams)
}
