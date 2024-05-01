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

	ctx2 := ctx.HideKind().Long()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `signature`).
		AddNonZero(ctx, `variadic`, sig.Variadic).
		AddNonZero(ctx2, `params`, sig.Params).
		AddNonZero(ctx2, `typeParams`, sig.TypeParams).
		AddNonZero(ctx.ShowKind().Short(), `return`, sig.Return)
}

func (sig *Signature) String() string {
	return jsonify.ToString(sig)
}
