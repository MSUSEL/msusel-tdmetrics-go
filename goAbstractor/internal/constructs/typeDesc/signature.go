package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature interface {
	TypeDesc
	_signature()
}

type SignatureArgs struct {
	RealType   *types.Signature
	Variadic   bool
	Params     []Named
	TypeParams []Named
	Return     TypeDesc
}

func NewSignature(reg Register, args SignatureArgs) Signature {
	return reg.RegisterSignature(&signatureImp{
		realType:   args.RealType,
		variadic:   args.Variadic,
		params:     args.Params,
		typeParams: args.TypeParams,
		returnType: args.Return,
	})
}

type signatureImp struct {
	realType *types.Signature

	variadic   bool
	params     []Named
	typeParams []Named
	returnType TypeDesc

	index int
}

func (sig *signatureImp) _signature() {}

func (sig *signatureImp) SetIndex(index int) {
	sig.index = index
}

func (sig *signatureImp) GoType() types.Type {
	return sig.realType
}

func (sig *signatureImp) Equal(other TypeDesc) bool {
	return equalTest(sig, other, func(a, b *signatureImp) bool {
		return a.variadic == b.variadic &&
			equal(a.returnType, b.returnType) &&
			equalList(a.params, b.params) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (sig *signatureImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, sig.index)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `signature`).
		AddNonZero(ctx, `variadic`, sig.variadic).
		AddNonZero(ctx.HideKind().Long(), `params`, sig.params).
		AddNonZero(ctx.HideKind().Long(), `typeParams`, sig.typeParams).
		AddNonZero(ctx.ShowKind().Short(), `return`, sig.returnType)
}

func (sig *signatureImp) String() string {
	return jsonify.ToString(sig)
}
