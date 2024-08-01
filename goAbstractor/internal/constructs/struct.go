package constructs

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
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
		Embedded []bool
		Package  Package
	}

	structImp struct {
		realType *types.Struct
		fields   []Named
		embedded []bool
		index    int
	}
)

func newStruct(args StructArgs) Struct {
	if len(args.Fields) != len(args.Embedded) {
		panic(terror.New(`the field count must be the same as the embedded flags`).
			With(`fields`, len(args.Fields)).
			With(`embedded`, len(args.Embedded)))
	}

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		pkg := args.Package.Source().Types
		fields := make([]*types.Var, len(args.Fields))
		for i, f := range args.Fields {
			fields[i] = types.NewField(token.NoPos, pkg, f.Name(), f.GoType(), args.Embedded[i])
		}

		args.RealType = types.NewStruct(fields, nil)
	}

	return &structImp{
		realType: args.RealType,
		fields:   args.Fields,
		embedded: args.Embedded,
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
		AddIf(ctx, ctx.IsIndexShown(), `index`, s.index).
		AddNonZero(ctx2, `fields`, s.fields)
}
