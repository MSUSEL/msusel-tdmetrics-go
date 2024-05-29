package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature interface {
	TypeDesc

	SetVariadic(v bool)
	SetReturn(t TypeDesc)
	AddParam(name string, t TypeDesc) Named
	AppendParam(tn ...Named)
	AppendTypeParam(tp ...Named)
}

type signatureImp struct {
	typ *types.Signature

	variadic   bool
	params     []Named
	typeParams []Named
	returnType TypeDesc

	index int
}

func NewSignature(typ *types.Signature) Signature {
	return &signatureImp{
		typ: typ,
	}
}

func (sig *signatureImp) SetIndex(index int) {
	sig.index = index
}

func (sig *signatureImp) GoType() types.Type {
	return sig.typ
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

func (sig *signatureImp) SetVariadic(v bool) {
	sig.variadic = v
}

func (sig *signatureImp) SetReturn(t TypeDesc) {
	sig.returnType = t
}

func (sig *signatureImp) AddParam(name string, t TypeDesc) Named {
	tn := NewNamed(name, t)
	sig.AppendParam(tn)
	return tn
}

func (sig *signatureImp) AppendParam(tn ...Named) {
	sig.params = append(sig.params, tn...)
}

func (sig *signatureImp) AppendTypeParam(tp ...Named) {
	sig.typeParams = append(sig.typeParams, tp...)
}
