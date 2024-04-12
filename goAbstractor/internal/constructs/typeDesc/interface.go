package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Interface struct {
	Index    int
	Inherits []*Interface
	Methods  []*Func
}

func (ti *Interface) _isTypeDesc() {}

func (ti *Interface) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.GetBool(`onlyIndex`) {
		return jsonify.New(ctx, ti.Index)
	}

	ctx = ctx.Copy().Set(`onlyIndex`, true)

	if len(ti.Methods) <= 0 {
		return jsonify.New(ctx, `object`)
	}

	return jsonify.NewMap().
		Add(ctx, `kind`, `interface`).
		AddNonZero(ctx, `inherits`, ti.Inherits).
		AddNonZero(ctx, `methods`, ti.Methods)
}

type Func struct {
	Name      string
	Signature *Signature
}

func (f *Func) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, f.Name).
		Add(ctx, `signature`, f.Signature)
}
