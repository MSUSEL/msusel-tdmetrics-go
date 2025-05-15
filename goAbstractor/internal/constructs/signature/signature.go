package signature

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type signatureImp struct {
	realType *types.Signature

	variadic bool
	params   []constructs.Argument
	results  []constructs.Argument

	index int
	alive bool
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

func (m *signatureImp) Kind() kind.Kind     { return kind.Signature }
func (m *signatureImp) Index() int          { return m.index }
func (m *signatureImp) SetIndex(index int)  { m.index = index }
func (m *signatureImp) Alive() bool         { return m.alive }
func (m *signatureImp) SetAlive(alive bool) { m.alive = alive }
func (m *signatureImp) GoType() types.Type  { return m.realType }

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
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.SliceComparerPend(aImp.params, bImp.params),
			constructs.SliceComparerPend(aImp.results, bImp.results),
			comp.DefaultPend(aImp.variadic, bImp.variadic),
		)
	}
}

func (m *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, m.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, m.Kind(), m.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, m.index).
		AddNonZero(ctx, `variadic`, m.variadic).
		AddNonZero(ctx.OnlyIndex(), `params`, m.params).
		AddNonZero(ctx.OnlyIndex(), `results`, m.results)
}

func (m *signatureImp) ToStringer(s stringer.Stringer) {
	s.Write(`func(`).
		WriteList(``, `, `, ``, m.params).
		Write(`)`)
	switch len(m.results) {
	case 0:
		break
	case 1:
		if len(m.results[0].Name()) > 0 {
			s.Write(`(`, m.results[0], `)`)
		} else {
			s.Write(` `, m.results[0])
		}
	default:
		s.WriteList(`(`, `, `, `)`, m.results)
	}
}

func (m *signatureImp) String() string {
	return stringer.String(m)
}
