package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Interface struct {
	Methods []*Func
}

func (ti *Interface) _isTypeDesc() {}

func (ti *Interface) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `interface`).
		AddNonZero(ctx, `methods`, ti.Methods)
}

type Func struct {
	Name      string
	Signature *Signature
}

func (f *Func) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, f.Name).
		Add(ctx, `signature`, f.Signature)
}
