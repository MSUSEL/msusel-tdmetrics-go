package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Method interface {
		Declaration
		_method()

		Metrics() metrics.Metrics
		Signature() Signature
		RecvName() string
		SetReceiver(recv ClassDecl)
		IsInit() bool
	}

	MethodArgs struct {
		Package    Package
		Name       string
		Location   locs.Loc
		Signature  Signature
		Metrics    metrics.Metrics
		NoCopyRecv bool
		RecvName   string
		TypeParams []Named
	}

	methodImp struct {
		pkg        Package
		name       string
		loc        locs.Loc
		signature  Signature
		metrics    metrics.Metrics
		noCopyRecv bool
		recvName   string
		receiver   ClassDecl
		typeParams []Named
		index      int
	}
)

func newMethod(args MethodArgs) Method {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`signature`, args.Signature)
	assert.ArgNotNil(`location`, args.Location)

	return &methodImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		signature:  args.Signature,
		metrics:    args.Metrics,
		noCopyRecv: args.NoCopyRecv,
		recvName:   args.RecvName,
		typeParams: args.TypeParams,
	}
}

func (m *methodImp) _method()                   {}
func (m *methodImp) Kind() kind.Kind            { return kind.Method }
func (m *methodImp) GoType() types.Type         { return m.signature.GoType() }
func (m *methodImp) SetIndex(index int)         { m.index = index }
func (m *methodImp) Name() string               { return m.name }
func (m *methodImp) Location() locs.Loc         { return m.loc }
func (m *methodImp) Package() Package           { return m.pkg }
func (m *methodImp) Metrics() metrics.Metrics   { return m.metrics }
func (m *methodImp) Signature() Signature       { return m.signature }
func (m *methodImp) RecvName() string           { return m.recvName }
func (m *methodImp) SetReceiver(recv ClassDecl) { m.receiver = recv }

func (m *methodImp) IsInit() bool {
	return strings.HasPrefix(m.name, `init#`) && m.signature.Vacant() && len(m.recvName) <= 0
}

func (m *methodImp) CompareTo(other Construct) int {
	b := other.(*methodImp)
	if cmp := Compare(m.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(m.name, b.name); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(m.recvName, b.recvName); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(m.typeParams, b.typeParams); cmp != 0 {
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
		AddIf(ctx, ctx.IsIndexShown(), `index`, m.index).
		Add(ctx2, `package`, m.pkg).
		Add(ctx2, `name`, m.name).
		Add(ctx2, `signature`, m.signature).
		AddNonZero(ctx2, `metrics`, m.metrics).
		AddNonZero(ctx2, `receiver`, m.receiver).
		AddNonZero(ctx2, `typeParams`, m.typeParams).
		AddNonZero(ctx2, `loc`, m.loc)

	if ctx.IsReceiverShown() {
		data.AddNonZero(ctx, `noCopyRecv`, m.noCopyRecv).
			AddNonZero(ctx, `recvName`, m.recvName)
	}
	return data
}

func (m *methodImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, m.signature)
}
