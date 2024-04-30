package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Method struct {
	Name      string
	Signature *typeDesc.Signature
	Metrics   metrics.Metrics

	NoCopyRecv bool
	Receiver   string
}

func NewMethod(name string, sig *typeDesc.Signature) *Method {
	return &Method{
		Name:      name,
		Signature: sig,
	}
}

func (m *Method) ToJson(ctx *jsonify.Context) jsonify.Datum {
	data := jsonify.NewMap().
		Add(ctx, `name`, m.Name).
		Add(ctx, `signature`, m.Signature).
		AddNonZero(ctx, `metrics`, m.Metrics)

	if ctx.ShowReceivers {
		data.AddNonZero(ctx, `noCopyRecv`, m.NoCopyRecv).
			AddNonZero(ctx, `receiver`, m.Receiver)
	}
	return data
}

func (m *Method) String() string {
	return jsonify.ToString(m)
}
