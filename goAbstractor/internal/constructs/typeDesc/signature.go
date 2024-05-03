package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Signature struct {
	Variadic   bool
	Params     []*Named
	TypeParams []*Named
	Return     TypeDesc

	Index int
}

func NewSignature() *Signature {
	return &Signature{}
}

func (sig *Signature) _isTypeDesc() {}

func (sig *Signature) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, sig.Index)
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
