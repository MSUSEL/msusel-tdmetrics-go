package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Signature struct {
	Index      int
	Variadic   bool
	Params     []*Param
	Return     TypeDesc
	TypeParams []*TypeParam
}

func (sig *Signature) _isTypeDesc() {}

func (sig *Signature) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.GetBool(`onlyIndex`) {
		return jsonify.New(ctx, sig.Index)
	}

	ctx2 := ctx.Copy().
		Remove(`noKind`).
		Set(`onlyIndex`, true)

	showKind := !ctx.GetBool(`noKind`)
	return jsonify.NewMap().
		AddIf(ctx2, showKind, `kind`, `signature`).
		AddNonZero(ctx2, `variadic`, sig.Variadic).
		AddNonZero(ctx2, `params`, sig.Params).
		AddNonZero(ctx2, `return`, sig.Return).
		AddNonZero(ctx2, `typeParams`, sig.TypeParams)
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
