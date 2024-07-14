package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Method interface {
		Definition
		_method()

		Signature() TypeDesc
		ReceiverName() string
	}

	MethodArgs struct {
		Package    Package
		Name       string
		Signature  TypeDesc
		Metrics    metrics.Metrics
		NoCopyRecv bool
		Receiver   string
	}

	methodImp struct {
		pkg        Package
		name       string
		signature  TypeDesc
		metrics    metrics.Metrics
		noCopyRecv bool
		recvName   string
		receiver   Class
		index      int
	}
)

func newMethod(args MethodArgs) Method {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`signature`, args.Signature)

	return &methodImp{
		pkg:        args.Package,
		name:       args.Name,
		signature:  args.Signature,
		metrics:    args.Metrics,
		noCopyRecv: args.NoCopyRecv,
		recvName:   args.Receiver,
	}
}

func (m *methodImp) _method()             {}
func (m *methodImp) Kind() kind.Kind      { return kind.Method }
func (m *methodImp) GoType() types.Type   { return m.signature.GoType() }
func (m *methodImp) SetIndex(index int)   { m.index = index }
func (m *methodImp) Name() string         { return m.name }
func (m *methodImp) Package() Package     { return m.pkg }
func (m *methodImp) Signature() TypeDesc  { return m.signature }
func (m *methodImp) ReceiverName() string { return m.recvName }

func (m *methodImp) CompareTo(other Construct) int {
	b := other.(*methodImp)
	if cmp := Compare(m.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(m.name, b.name); cmp != 0 {
		return cmp
	}
	return Compare(m.signature, b.signature)
}

func (m *methodImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, m.index)
	}

	ctx2 := ctx.HideKind().Short()
	data := jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, m.Kind()).
		Add(ctx2, `package`, m.pkg).
		Add(ctx2, `name`, m.name).
		Add(ctx2, `signature`, m.signature).
		AddNonZero(ctx2, `metrics`, m.metrics).
		AddNonZero(ctx2, `receiver`, m.receiver)

	if ctx.IsReceiverShown() {
		data.AddNonZero(ctx, `noCopyRecv`, m.noCopyRecv).
			AddNonZero(ctx, `recvName`, m.recvName)
	}
	return data
}

func (m *methodImp) Visit(v visitor.Visitor) bool {
	return visitor.Visit(v, m.signature)
}
