package typeDesc

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/wrapKind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Wrap struct {
	Kind wrapKind.WrapKind
	Elem TypeDesc
}

func (tw *Wrap) _isTypeDesc() {}

func (tw *Wrap) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `kind`, string(tw.Kind)).
		Add(ctx, `elem`, tw.Elem)
}
