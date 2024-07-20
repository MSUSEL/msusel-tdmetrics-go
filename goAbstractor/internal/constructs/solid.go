package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	// Solid represents a generic type that has been resolved to a specific type
	// with specific type parameters, e.g. List<T> might be resolved to List<int>.
	// The type parameter resolution may be referencing another type parameter,
	// e.g. a method signature inside a generic interface.
	Solid interface {
		TypeDesc
		_solid()
	}

	SolidArgs struct {
		RealType   types.Type
		Target     TypeDesc
		TypeParams []TypeDesc
	}

	solidImp struct {
		realType types.Type

		index      int
		target     TypeDesc
		typeParams []TypeDesc
	}
)

func newSolid(args SolidArgs) Solid {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`target`, args.Target)
	assert.ArgNotEmpty(`type params`, args.TypeParams)

	return &solidImp{
		realType:   args.RealType,
		target:     args.Target,
		typeParams: args.TypeParams,
	}
}

func (s *solidImp) _solid()            {}
func (s *solidImp) Kind() kind.Kind    { return kind.Solid }
func (s *solidImp) SetIndex(index int) { s.index = index }
func (s *solidImp) GoType() types.Type { return s.realType }

func (s *solidImp) CompareTo(other Construct) int {
	b := other.(*solidImp)
	if cmp := Compare(s.target, b.target); cmp != 0 {
		return cmp
	}
	return CompareSlice(s.typeParams, b.typeParams)
}

func (s *solidImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, s.target)
	visitor.Visit(v, s.typeParams...)
}

func (s *solidImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, s.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, s.Kind()).
		AddNonZero(ctx2, `target`, s.target).
		AddNonZero(ctx2, `typeParams`, s.typeParams)
}
