package method

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type methodImp struct {
	realType *types.Signature
	pkg      constructs.Package
	name     string
	loc      locs.Loc

	typeParams []constructs.TypeParam
	signature  constructs.Signature
	metrics    constructs.Metrics
	recvName   string
	receiver   constructs.Object
	noCopyRecv bool

	instances collections.SortedSet[constructs.Instance]

	id any
}

func newMethod(args constructs.MethodArgs) constructs.Method {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgNotNil(`signature`, args.Signature)
	assert.ArgNotNil(`type params`, args.TypeParams)
	assert.ArgNotNil(`location`, args.Location)

	if !utils.IsNil(args.Receiver) {
		rName := args.Receiver.Name()
		if len(args.RecvName) > 0 && args.RecvName != rName {
			panic(terror.New(`name of receiver and a receiver were both given but the name didn't match the receiver`).
				With(`receiver name`, args.RecvName).
				With(`receiver`, rName))
		}
		args.RecvName = rName
	}

	met := &methodImp{
		pkg:  args.Package,
		name: args.Name,
		loc:  args.Location,

		typeParams: args.TypeParams,
		signature:  args.Signature,
		metrics:    args.Metrics,
		recvName:   args.RecvName,
		receiver:   args.Receiver,
		noCopyRecv: args.NoCopyRecv,

		instances: sortedSet.New(instance.Comparer()),
	}

	if !utils.IsNil(met.receiver) {
		return met.receiver.AddMethod(met)
	}
	return met
}

func (m *methodImp) IsDeclaration() {}
func (m *methodImp) IsMethod()      {}

func (m *methodImp) Kind() kind.Kind    { return kind.Method }
func (m *methodImp) Id() any            { return m.id }
func (m *methodImp) SetId(id any)       { m.id = id }
func (m *methodImp) GoType() types.Type { return m.realType }

func (m *methodImp) Package() constructs.Package { return m.pkg }
func (m *methodImp) Name() string                { return m.name }
func (m *methodImp) Location() locs.Loc          { return m.loc }

func (m *methodImp) Type() constructs.TypeDesc          { return m.signature }
func (m *methodImp) Signature() constructs.Signature    { return m.signature }
func (m *methodImp) Metrics() constructs.Metrics        { return m.metrics }
func (m *methodImp) TypeParams() []constructs.TypeParam { return m.typeParams }

func (m *methodImp) Instances() collections.ReadonlySortedSet[constructs.Instance] {
	return m.instances.Readonly()
}

func (m *methodImp) AddInstance(inst constructs.Instance) constructs.Instance {
	v, _ := m.instances.TryAdd(inst)
	return v
}

func (m *methodImp) ReceiverName() string               { return m.recvName }
func (m *methodImp) SetReceiver(recv constructs.Object) { m.receiver = recv }
func (m *methodImp) Receiver() constructs.Object        { return m.receiver }
func (m *methodImp) NoCopyRecv() bool                   { return m.noCopyRecv }

func (m *methodImp) NeedsReceiver() bool {
	return utils.IsNil(m.receiver) && len(m.recvName) > 0
}

func (m *methodImp) IsInit() bool {
	return strings.HasPrefix(m.name, `init#`) &&
		len(m.recvName) <= 0 &&
		len(m.typeParams) <= 0 &&
		m.signature.IsVacant()
}

func (m *methodImp) IsNamed() bool {
	return len(m.name) > 0
}

func (m *methodImp) IsGeneric() bool {
	return len(m.typeParams) > 0
}

func (m *methodImp) HasReceiver() bool {
	return !utils.IsNil(m.receiver) || len(m.recvName) > 0
}

func (m *methodImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Method](m, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Method] {
	return func(a, b constructs.Method) int {
		aImp, bImp := a.(*methodImp), b.(*methodImp)
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			comp.DefaultPend(aImp.recvName, bImp.recvName),
			constructs.SliceComparerPend(aImp.typeParams, bImp.typeParams),
			constructs.ComparerPend(aImp.signature, bImp.signature),
		)
	}
}

func (m *methodImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, m.id)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsIdShown(), `id`, m.id).
		Add(ctx2, `package`, m.pkg).
		Add(ctx2, `name`, m.name).
		AddNonZero(ctx2, `loc`, m.loc).
		AddNonZero(ctx2, `typeParams`, m.typeParams).
		Add(ctx2, `signature`, m.signature).
		AddNonZero(ctx2, `metrics`, m.metrics).
		AddNonZero(ctx2, `instances`, m.instances.ToSlice()).
		AddNonZero(ctx2, `receiver`, m.receiver).
		AddNonZeroIf(ctx2, ctx.IsReceiverShown(), `recvName`, m.recvName)
}

func (m *methodImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(m.name + ` `)
	if len(m.typeParams) > 0 {
		buf.WriteString(`[` + m.instances.Enumerate().Join(`, `) + `]`)
	}
	buf.WriteString(m.signature.String())
	return buf.String()
}
