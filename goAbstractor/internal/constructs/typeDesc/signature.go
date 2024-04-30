package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Signature struct {
	Variadic   bool
	Params     []*Param
	TypeParams []*Param
	Return     TypeDesc

	Index int
}

func NewSignature() *Signature {
	return &Signature{}
}

func (sig *Signature) _isTypeDesc() {}

func (sig *Signature) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.Short {
		return jsonify.New(ctx, sig.Index)
	}

	ctx2 := ctx.Copy()
	ctx2.NoKind = false
	ctx2.Short = true

	return jsonify.NewMap().
		AddIf(ctx2, !ctx.NoKind, `kind`, `signature`).
		AddNonZero(ctx2, `variadic`, sig.Variadic).
		AddNonZero(ctx2, `params`, sig.Params).
		AddNonZero(ctx2, `return`, sig.Return)
}

func (sig *Signature) String() string {
	return jsonify.ToString(sig)
}

type Param struct {
	Name string
	Type TypeDesc
}

func (p *Param) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `name`, p.Name).
		Add(ctx, `type`, p.Type)
}

func (p *Param) String() string {
	return jsonify.ToString(p)
}
