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
	Signature interface {
		TypeDesc
		Vacant() bool
		_signature()
	}

	SignatureArgs struct {
		RealType   *types.Signature
		Variadic   bool
		Params     []Named
		TypeParams []Named
		Return     TypeDesc

		// Package is only needed if the real type is nil
		// so that a Go signature type has to be created.
		Package Package
	}

	signatureImp struct {
		realType *types.Signature

		variadic   bool
		params     []Named
		typeParams []Named
		returnType TypeDesc

		index int
	}
)

func newSignature(args SignatureArgs) Signature {
	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		pkg := args.Package.Source().Types
		tp := make([]*types.TypeParam, len(args.TypeParams))
		for i, t := range args.TypeParams {
			tName := types.NewTypeName(token.NoPos, pkg, ``, t.GoType())
			tp[i] = types.NewTypeParam(tName, t.GoType())
		}

		params := make([]*types.Var, len(args.Params))
		for i, p := range args.Params {
			params[i] = types.NewVar(token.NoPos, pkg, p.Name(), p.GoType())
		}

		var ret *types.Tuple
		if !utils.IsNil(args.Return) {
			v := types.NewVar(token.NoPos, pkg, ``, args.Return.GoType())
			ret = types.NewTuple(v)
		}

		args.RealType = types.NewSignatureType(nil, nil,
			tp, types.NewTuple(params...), ret, args.Variadic)
	}

	return &signatureImp{
		realType:   args.RealType,
		variadic:   args.Variadic,
		params:     args.Params,
		typeParams: args.TypeParams,
		returnType: args.Return,
	}
}

func (s *signatureImp) _signature()        {}
func (s *signatureImp) Kind() kind.Kind    { return kind.Signature }
func (s *signatureImp) SetIndex(index int) { s.index = index }
func (s *signatureImp) GoType() types.Type { return s.realType }

func (s *signatureImp) Vacant() bool {
	return len(s.params) <= 0 &&
		len(s.typeParams) <= 0 &&
		utils.IsNil(s.returnType)
}

func (s *signatureImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, s.params...)
	visitor.Visit(v, s.typeParams...)
	visitor.Visit(v, s.returnType)
}

func (s *signatureImp) CompareTo(other Construct) int {
	b := other.(*signatureImp)
	if !s.variadic && b.variadic {
		return -1
	}
	if s.variadic && !b.variadic {
		return 1
	}
	if cmp := CompareSlice(s.typeParams, b.typeParams); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(s.params, b.params); cmp != 0 {
		return cmp
	}
	return Compare(s.returnType, b.returnType)
}

func (s *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, s.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, s.Kind()).
		AddNonZero(ctx2, `variadic`, s.variadic).
		AddNonZero(ctx2, `params`, s.params).
		AddNonZero(ctx2, `typeParams`, s.typeParams).
		AddNonZero(ctx2, `return`, s.returnType)
}
