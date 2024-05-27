package typeDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Signature struct {
	typ *types.Signature

	Variadic   bool
	Params     []*Named
	TypeParams []*Named
	Return     TypeDesc

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
		return a.Variadic == b.Variadic &&
			equal(a.Return, b.Return) &&
			equalList(a.Params, b.Params) &&
			equalList(a.TypeParams, b.TypeParams)
	})
}

func (sig *Signature) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, sig.index)
	}

	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `signature`).
		AddNonZero(ctx, `variadic`, sig.Variadic).
		AddNonZero(ctx.HideKind().Long(), `params`, sig.Params).
		AddNonZero(ctx.HideKind().Long(), `typeParams`, sig.TypeParams).
		AddNonZero(ctx.ShowKind().Short(), `return`, sig.Return)
}

func (sig *Signature) String() string {
	return jsonify.ToString(sig)
}

func (sig *Signature) AddParam(name string, t TypeDesc) *Named {
	tn := NewNamed(name, t)
	sig.Params = append(sig.Params, tn)
	return tn
}

func (sig *Signature) AppendTypeParam(tp ...*Named) {
	sig.TypeParams = append(sig.TypeParams, tp...)
}
