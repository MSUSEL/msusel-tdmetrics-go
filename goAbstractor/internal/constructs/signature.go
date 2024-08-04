package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature interface {
	TypeDesc
	Vacant() bool
	_signature()
}

type SignatureArgs struct {
	RealType *types.Signature
	Variadic bool
	Params   []TypeDesc
	Returns  []TypeDesc
}

type signatureImp struct {
	realType *types.Signature
	variadic bool
	params   []TypeDesc
	returns  []TypeDesc
	index    int
}

func newSignature(args SignatureArgs) Signature {
	assert.ArgNotNil(`real type`, args.RealType)

	return &signatureImp{
		realType: args.RealType,
		variadic: args.Variadic,
		params:   args.Params,
		returns:  args.Returns,
	}
}

func (s *signatureImp) _signature()        {}
func (s *signatureImp) Kind() kind.Kind    { return kind.Signature }
func (s *signatureImp) setIndex(index int) { s.index = index }
func (s *signatureImp) GoType() types.Type { return s.realType }

func (s *signatureImp) Vacant() bool {
	return len(s.params) <= 0 && len(s.returns) <= 0
}

func (s *signatureImp) compareTo(other Construct) int {
	b := other.(*signatureImp)
	return or(
		func() int { return boolCompare(s.variadic, b.variadic) },
		func() int { return compareSlice(s.params, b.params) },
		func() int { return compareSlice(s.returns, b.returns) },
	)
}

func (s *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, s.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, s.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, s.index).
		AddNonZero(ctx2, `variadic`, s.variadic).
		AddNonZero(ctx2, `params`, s.params).
		AddNonZero(ctx2, `returns`, s.returns)
}
