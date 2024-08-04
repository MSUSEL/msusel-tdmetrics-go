package constructs

import (
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Method interface {
	Construct
	_method()

	Package() Package
	Name() string
	Location() locs.Loc

	Signature() Signature
	Metrics() metrics.Metrics

	addInstance(inst Instance) Instance
	receiverName() string
	setReceiver(recv Declaration)

	IsInit() bool
	IsNamed() bool
	IsGeneric() bool
}

type MethodArgs struct {
	Package  Package
	Name     string
	Location locs.Loc

	Signature  Signature
	TypeParams []TypeParam

	Metrics metrics.Metrics

	NoCopyRecv bool
	RecvName   string
}

type methodImp struct {
	pkg  Package
	name string
	loc  locs.Loc

	signature  Signature
	typeParams []TypeParam
	instances  Set[Instance]

	metrics metrics.Metrics

	noCopyRecv bool
	recvName   string
	receiver   Declaration

	index int
}

func newMethod(args MethodArgs) Method {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`signature`, args.Signature)
	assert.ArgNotNil(`type params`, args.TypeParams)

	return &methodImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		signature:  args.Signature,
		metrics:    args.Metrics,
		noCopyRecv: args.NoCopyRecv,
		recvName:   args.RecvName,
		typeParams: args.TypeParams,

		instances: NewSet[Instance](),
	}
}

func (m *methodImp) _method()           {}
func (m *methodImp) Kind() kind.Kind    { return kind.Method }
func (m *methodImp) setIndex(index int) { m.index = index }

func (m *methodImp) Package() Package   { return m.pkg }
func (m *methodImp) Name() string       { return m.name }
func (m *methodImp) Location() locs.Loc { return m.loc }

func (m *methodImp) Metrics() metrics.Metrics { return m.metrics }
func (m *methodImp) Signature() Signature     { return m.signature }

func (m *methodImp) addInstance(inst Instance) Instance {
	return m.instances.Insert(inst)
}

func (m *methodImp) receiverName() string         { return m.recvName }
func (m *methodImp) setReceiver(recv Declaration) { m.receiver = recv }

func (m *methodImp) IsInit() bool {
	return strings.HasPrefix(m.name, `init#`) &&
		m.signature.Vacant() &&
		len(m.recvName) <= 0
}

func (m *methodImp) IsNamed() bool {
	return len(m.name) > 0
}

func (m *methodImp) IsGeneric() bool {
	return len(m.typeParams) > 0
}

func (m *methodImp) compareTo(other Construct) int {
	b := other.(*methodImp)
	return or(
		func() int { return Compare(m.pkg, b.pkg) },
		func() int { return strings.Compare(m.name, b.name) },
		func() int { return strings.Compare(m.recvName, b.recvName) },
		func() int { return compareSlice(m.typeParams, b.typeParams) },
		func() int { return Compare(m.signature, b.signature) },
	)
}

func (m *methodImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, m.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, m.index).
		AddNonZero(ctx2, `package`, m.pkg).
		AddNonZero(ctx2, `name`, m.name).
		AddNonZero(ctx2, `loc`, m.loc).
		AddNonZero(ctx2, `signature`, m.signature).
		AddNonZero(ctx2, `metrics`, m.metrics).
		AddNonZero(ctx2, `typeParams`, m.typeParams).
		AddNonZero(ctx2, `instances`, m.instances).
		AddNonZero(ctx2, `receiver`, m.receiver).
		AddNonZeroIf(ctx2, ctx.IsReceiverShown(), `noCopyRecv`, m.noCopyRecv).
		AddNonZeroIf(ctx2, ctx.IsReceiverShown(), `recvName`, m.recvName)
}
