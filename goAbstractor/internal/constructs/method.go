package constructs

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type (
	Method interface {
		Construct
		_method()

		Name() string
		Signature() TypeDesc
		Receiver() string
	}

	MethodArgs struct {
		Name       string
		Signature  TypeDesc
		Metrics    metrics.Metrics
		NoCopyRecv bool
		Receiver   string
	}

	methodImp struct {
		name       string
		signature  TypeDesc
		metrics    metrics.Metrics
		noCopyRecv bool
		receiver   string
	}
)

func newMethod(args MethodArgs) Method {
	return &methodImp{
		name:       args.Name,
		signature:  args.Signature,
		metrics:    args.Metrics,
		noCopyRecv: args.NoCopyRecv,
		receiver:   args.Receiver,
	}
}

func (m *methodImp) _method() {}

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
