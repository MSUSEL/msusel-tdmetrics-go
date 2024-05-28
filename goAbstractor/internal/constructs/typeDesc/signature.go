package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature struct {
	typ *types.Signature

	variadic   bool
	params     []Named
	typeParams []Named
	returnType TypeDesc

	index int
}

func NewSignature(typ *types.Signature) *Signature {
	return &Signature{
		typ: typ,
	}
}

func (sig *Signature) SetIndex(index int) {
	sig.index = index
}

func (sig *Signature) GoType() types.Type {
	return sig.typ
}

func (sig *Signature) Equal(other TypeDesc) bool {
	return equalTest(sig, other, func(a, b *Signature) bool {
		return a.variadic == b.variadic &&
			equal(a.returnType, b.returnType) &&
			equalList(a.params, b.params) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (sig *Signature) ToJson(ctx *jsonify.Context) jsonify.Datum {
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

func (sig *Signature) String() string {
	return jsonify.ToString(sig)
}

func (sig *Signature) SetVariadic(v bool) {
	sig.variadic = v
}

func (sig *Signature) SetReturn(t TypeDesc) {
	sig.returnType = t
}

func (sig *Signature) AddParam(name string, t TypeDesc) Named {
	tn := NewNamed(name, t)
	sig.AppendParam(tn)
	return tn
}

func (sig *Signature) AppendParam(tn ...Named) {
	sig.params = append(sig.params, tn...)
}

func (sig *Signature) AppendTypeParam(tp ...Named) {
	sig.typeParams = append(sig.typeParams, tp...)
}
