package constructs

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Struct interface {
		TypeDesc
		_struct()
	}

	StructArgs struct {
		RealType *types.Struct
		Fields   []Named

		// Package is only needed if the real type is nil
		// so that a Go interface type has to be created.
		Package Package
	}

	structImp struct {
		realType *types.Struct
		fields   []Named
		index    int
	}
)

func newStruct(args StructArgs) Struct {
	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		fields := make([]*types.Var, len(args.Fields))
		pkg := args.Package.Source().Types
		for i, f := range args.Fields {
			fields[i] = types.NewVar(token.NoPos, pkg, f.Name(), f.GoType())
		}

		args.RealType = types.NewStruct(fields, nil)
	}

	return &structImp{
		realType: args.RealType,
		fields:   args.Fields,
	}
}

func (s *structImp) _struct()           {}
func (s *structImp) Kind() kind.Kind    { return kind.Struct }
func (s *structImp) SetIndex(index int) { s.index = index }
func (s *structImp) GoType() types.Type { return s.realType }

func (s *structImp) CompareTo(other Construct) int {
	return CompareSlice(s.fields, other.(*structImp).fields)
}

func (s *structImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, s.fields...)
}

func (s *structImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, s.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, s.Kind()).
		AddNonZero(ctx2, `fields`, s.fields)
}
