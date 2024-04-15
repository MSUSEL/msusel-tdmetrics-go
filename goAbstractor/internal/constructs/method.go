package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Method struct {
	Name      string
	Signature *typeDesc.Signature

	NoCopyRecv bool
	Receiver   string
}

func (m *Method) ToJson(ctx *jsonify.Context) jsonify.Datum {
	data := jsonify.NewMap().
		Add(ctx, `name`, m.Name).
		Add(ctx, `signature`, m.Signature)

	if ctx.GetBool(`showReceivers`) {
		data.AddNonZero(ctx, `noCopyRecv`, m.NoCopyRecv).
			AddNonZero(ctx, `receiver`, m.Receiver)
	}
	return data
}

func (m *Method) String() string {
	return jsonify.ToString(m)
}