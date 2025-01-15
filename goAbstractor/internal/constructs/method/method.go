package method

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/methodInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type methodImp struct {
	realType *types.Signature
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc
	index    int
	alive    bool

	typeParams []constructs.TypeParam
	signature  constructs.Signature
	metrics    constructs.Metrics
	recvName   string
	receiver   constructs.Object
	ptrRecv    bool

	instances collections.SortedSet[constructs.MethodInst]
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
		pkg:        args.Package,
		name:       args.Name,
		exported:   args.Exported,
		loc:        args.Location,
		typeParams: args.TypeParams,
		signature:  args.Signature,
		metrics:    args.Metrics,
		recvName:   args.RecvName,
		receiver:   args.Receiver,
		ptrRecv:    args.PointerRecv,
		instances:  sortedSet.New(methodInst.Comparer()),
	}

	if !utils.IsNil(met.receiver) {
		return met.receiver.AddMethod(met)
	}
	return met
}

func (m *methodImp) IsDeclaration() {}
func (m *methodImp) IsMethod()      {}

func (m *methodImp) Kind() kind.Kind     { return kind.Method }
func (m *methodImp) Index() int          { return m.index }
func (m *methodImp) SetIndex(index int)  { m.index = index }
func (m *methodImp) Alive() bool         { return m.alive }
func (m *methodImp) SetAlive(alive bool) { m.alive = alive }
func (m *methodImp) GoType() types.Type  { return m.realType }
func (m *methodImp) Name() string        { return m.name }
func (m *methodImp) Exported() bool      { return m.exported }
func (m *methodImp) Location() locs.Loc  { return m.loc }

func (m *methodImp) Package() constructs.Package        { return m.pkg }
func (m *methodImp) Type() constructs.TypeDesc          { return m.signature }
func (m *methodImp) Signature() constructs.Signature    { return m.signature }
func (m *methodImp) Metrics() constructs.Metrics        { return m.metrics }
func (m *methodImp) TypeParams() []constructs.TypeParam { return m.typeParams }

func (m *methodImp) Instances() collections.ReadonlySortedSet[constructs.MethodInst] {
	return m.instances.Readonly()
}

func (m *methodImp) AddInstance(inst constructs.MethodInst) constructs.MethodInst {
	v, _ := m.instances.TryAdd(inst)
	return v
}

func (m *methodImp) FindInstance(instanceTypes []constructs.TypeDesc) (constructs.MethodInst, bool) {
	cmp := constructs.SliceComparer[constructs.TypeDesc]()
	return m.instances.Enumerate().Where(func(i constructs.MethodInst) bool {
		return cmp(instanceTypes, i.InstanceTypes()) == 0
	}).First()
}

func (m *methodImp) ReceiverName() string               { return m.recvName }
func (m *methodImp) SetReceiver(recv constructs.Object) { m.receiver = recv }
func (m *methodImp) Receiver() constructs.Object        { return m.receiver }
func (m *methodImp) PointerRecv() bool                  { return m.ptrRecv }

func (m *methodImp) NeedsReceiver() bool {
	return utils.IsNil(m.receiver) && len(m.recvName) > 0
}

func (m *methodImp) IsInit() bool {
	return strings.HasPrefix(m.name, `init#`) &&
		m.HasReceiver() &&
		m.IsGeneric() &&
		m.signature.IsVacant()
}

func (m *methodImp) IsMain() bool {
	return m.name == `main` &&
		m.pkg.Name() == `main` &&
		m.HasReceiver() &&
		m.IsGeneric() &&
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
		if aImp == bImp {
			return 0
		}
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
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, m.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, m.Kind(), m.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, m.index).
		Add(ctx.OnlyIndex(), `package`, m.pkg).
		Add(ctx, `name`, m.name).
		AddNonZero(ctx, `loc`, m.loc).
		AddNonZero(ctx, `exported`, m.exported).
		AddNonZero(ctx.OnlyIndex(), `typeParams`, m.typeParams).
		Add(ctx.OnlyIndex(), `signature`, m.signature).
		AddNonZero(ctx.OnlyIndex(), `metrics`, m.metrics).
		AddNonZero(ctx.OnlyIndex(), `instances`, m.instances.ToSlice()).
		AddNonZero(ctx.OnlyIndex(), `receiver`, m.receiver).
		AddNonZero(ctx, `ptrRecv`, m.ptrRecv).
		AddNonZeroIf(ctx, ctx.IsDebugReceiverIncluded(), `recvName`, m.recvName)
}

func (m *methodImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`func `)
	buf.WriteString(m.pkg.Path())
	buf.WriteString(`.`)
	if len(m.recvName) > 0 {
		buf.WriteString(m.recvName)
		buf.WriteString(`.`)
	}
	buf.WriteString(m.name)
	if len(m.typeParams) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(m.typeParams...).Join(`, `))
		buf.WriteString(`]`)
	}
	sig := m.signature.String()
	sig, _ = strings.CutPrefix(sig, `func`)
	buf.WriteString(sig)
	return buf.String()
}
