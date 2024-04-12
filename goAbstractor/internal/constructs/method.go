package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Method struct {
	Name      string
	Signature *typeDesc.Signature
	Receiver  typeDesc.TypeDesc
}

func (m *Method) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `name`, m.Name).
		Add(ctx, `signature`, m.Signature).
		AddNonZero(ctx, `receiver`, m.Receiver)
}
