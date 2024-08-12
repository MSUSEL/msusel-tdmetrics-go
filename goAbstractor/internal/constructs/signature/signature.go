package signature

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type signatureImp struct {
	realType *types.Signature

	variadic bool
	params   []constructs.Argument
	results  []constructs.Argument

	index int
}

func createTuple(pkg *types.Package, args []constructs.Argument) *types.Tuple {
	vars := make([]*types.Var, len(args))
	for i, p := range args {
		vars[i] = types.NewVar(token.NoPos, pkg, p.Name(), p.Type().GoType())
	}
	return types.NewTuple(vars...)
}

func newSignature(args constructs.SignatureArgs) constructs.Signature {
	assert.ArgNoNils(`params`, args.Params)
	assert.ArgNoNils(`results`, args.Results)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)
		pkg := args.Package
		params := createTuple(pkg, args.Params)
		results := createTuple(pkg, args.Results)
		args.RealType = types.NewSignatureType(nil, nil, nil, params, results, args.Variadic)
	}

	return &signatureImp{
		realType: args.RealType,
		variadic: args.Variadic,
		params:   args.Params,
		results:  args.Results,
	}
}

func (m *signatureImp) IsTypeDesc()        {}
func (m *signatureImp) IsSignature()       {}
func (m *signatureImp) Kind() kind.Kind    { return kind.Signature }
func (m *signatureImp) SetIndex(index int) { m.index = index }
func (m *signatureImp) GoType() types.Type { return m.realType }

func (m *signatureImp) IsVacant() bool {
	return len(m.params) <= 0 && len(m.results) <= 0
}

func (s *signatureImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Signature](s, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Signature] {
	return func(a, b constructs.Signature) int {
		aImp, bImp := a.(*signatureImp), b.(*signatureImp)
		return comp.Or(
			constructs.SliceComparerPend(aImp.params, bImp.params),
			constructs.SliceComparerPend(aImp.results, bImp.results),
			comp.DefaultPend(aImp.variadic, bImp.variadic),
		)
	}
}

func (m *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, m.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, m.index).
		AddNonZero(ctx2, `variadic`, m.variadic).
		AddNonZero(ctx2, `params`, m.params).
		AddNonZero(ctx2, `results`, m.results)
}
