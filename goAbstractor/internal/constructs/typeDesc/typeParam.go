package typeDesc

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

// TODO: Need to rework
// Type param is defined on each param/return and signature right now.
// At minimum the params/returns could be just index references.
// Need to rework to use minimum common interfaces to be like Java.
// This means things like `int` need to have a pseudo interface.

type TypeParam struct {
	Index      int
	Constraint TypeDesc
}

func (tp *TypeParam) _isTypeDesc() {}

func (tp *TypeParam) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, `typeParam`).
		Add(ctx, `index`, tp.Index).
		Add(ctx, `constraint`, tp.Constraint)
}
