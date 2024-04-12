package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Signature struct {
	Variadic   bool
	Params     []*Param
	Return     TypeDesc
	TypeParams []*TypeParam
}

func (sig *Signature) _isTypeDesc() {}

func (sig *Signature) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `signature`).
		AddNonZero(ctx, `variadic`, sig.Variadic).
		AddNonZero(ctx, `params`, sig.Params).
		AddNonZero(ctx, `return`, sig.Return).
		AddNonZero(ctx, `typeParams`, sig.TypeParams)
}

type Param struct {
	Name string
	Type TypeDesc
}

func (p *Param) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `name`, p.Name).
		Add(ctx, `type`, p.Type)
}
