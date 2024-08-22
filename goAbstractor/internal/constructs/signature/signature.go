package signature

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

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

func createTuple(pkg *packages.Package, args []constructs.Argument) *types.Tuple {
	vars := make([]*types.Var, len(args))
	for i, p := range args {
		vars[i] = types.NewVar(token.NoPos, pkg.Types, p.Name(), p.Type().GoType())
	}
	return types.NewTuple(vars...)
}

func newSignature(args constructs.SignatureArgs) constructs.Signature {
	assert.ArgHasNoNils(`params`, args.Params)
	assert.ArgHasNoNils(`results`, args.Results)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)
		params := createTuple(args.Package, args.Params)
		results := createTuple(args.Package, args.Results)
		args.RealType = types.NewSignatureType(nil, nil, nil, params, results, args.Variadic)
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &signatureImp{
		realType: args.RealType,
		variadic: args.Variadic,
		params:   args.Params,
		results:  args.Results,
	}
}

func (m *signatureImp) IsTypeDesc()  {}
func (m *signatureImp) IsSignature() {}

func (m *signatureImp) Kind() kind.Kind    { return kind.Signature }
func (m *signatureImp) SetIndex(index int) { m.index = index }
func (m *signatureImp) GoType() types.Type { return m.realType }

func (m *signatureImp) Variadic() bool                 { return m.variadic }
func (m *signatureImp) Params() []constructs.Argument  { return m.params }
func (m *signatureImp) Results() []constructs.Argument { return m.results }

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

func (m *signatureImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`func(`)
	buf.WriteString(enumerator.Enumerate(m.params).Join(`, `))
	buf.WriteString(`)`)
	switch len(m.results) {
	case 0:
		break
	case 1:
		buf.WriteString(` ` + m.results[0].String())
	default:
		buf.WriteString(`(` + enumerator.Enumerate(m.results).Join(`, `) + `)`)
	}
	return buf.String()
}
