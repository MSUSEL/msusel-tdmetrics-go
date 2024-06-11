package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

// TODO: Need to handle when the type parameters are named differently
// in the receiver from the ones defined on the receiver's type def.

type Method interface {
	Visitable
	Name() string
	Signature() TypeDesc
	Receiver() string
	SetMetrics(met metrics.Metrics)
}

type methodImp struct {
	name       string
	signature  TypeDesc
	metrics    metrics.Metrics
	noCopyRecv bool
	receiver   string
}

func NewMethod(name string, sig TypeDesc, noCopyRecv bool, receiver string) Method {
	return &methodImp{
		name:       name,
		signature:  sig,
		noCopyRecv: noCopyRecv,
		receiver:   receiver,
	}
}

func (m *methodImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	data := jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `method`).
		Add(ctx2, `name`, m.name).
		Add(ctx2, `signature`, m.signature).
		AddNonZero(ctx2, `metrics`, m.metrics)

	if ctx.IsReceiverShown() {
		data.AddNonZero(ctx, `noCopyRecv`, m.noCopyRecv).
			AddNonZero(ctx, `receiver`, m.receiver)
	}
	return data
}

func (m *methodImp) Visit(v Visitor) {
	visitTest(v, m.signature)
}

func (m *methodImp) String() string {
	return jsonify.ToString(m)
}

func (m *methodImp) Name() string {
	return m.name
}

func (m *methodImp) Signature() TypeDesc {
	return m.signature
}

func (m *methodImp) Receiver() string {
	return m.receiver
}

func (m *methodImp) SetMetrics(met metrics.Metrics) {
	m.metrics = met
}
