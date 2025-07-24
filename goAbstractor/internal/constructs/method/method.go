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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/methodInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type methodImp struct {
	constructs.ConstructCore
	funcType *types.Func
	sigType  *types.Signature
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc

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
		funcType:   args.FuncType,
		sigType:    args.SigType,
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

func (m *methodImp) Kind() kind.Kind    { return kind.Method }
func (m *methodImp) GoType() types.Type { return m.sigType }
func (m *methodImp) Name() string       { return m.name }
func (m *methodImp) Exported() bool     { return m.exported }
func (m *methodImp) Location() locs.Loc { return m.loc }

func (m *methodImp) FuncType() *types.Func              { return m.funcType }
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
		m.IsConcreteFunc() &&
		m.signature.IsVacant()
}

func (m *methodImp) IsMain() bool {
	return m.name == `main` &&
		m.pkg.Name() == `main` &&
		m.IsConcreteFunc() &&
		m.signature.IsVacant()
}

func (m *methodImp) IsTester() bool {
	return m.IsTest() ||
		m.IsBenchmark() ||
		m.IsFuzz() ||
		m.IsExample()
}

// hasTestParam determines if this method follows the pattern
// `func(*testing.<name>)` such as `func(*testing.T)`.
func (m *methodImp) hasTestParam(name string) bool {
	s := m.signature
	if utils.IsNil(s) || len(s.Results()) != 0 || len(s.Params()) != 1 {
		return false
	}
	p := s.Params()[0]
	if utils.IsNil(p) {
		return false
	}
	t := p.Type()
	if utils.IsNil(t) || t.Kind() != kind.InterfaceInst {
		return false
	}
	i, ok := t.(constructs.InterfaceInst)
	if !ok || i.Resolved().Hint() != hint.Pointer || len(i.ImplicitTypes()) != 1 {
		return false
	}
	ta := i.ImplicitTypes()[0]
	if utils.IsNil(ta) || ta.Kind() != kind.Object {
		return false
	}
	o, ok := ta.(constructs.Object)
	return ok &&
		!utils.IsNil(o.Package()) &&
		o.Package().Path() == `testing` &&
		o.Name() == name
}

func (m *methodImp) IsTest() bool {
	return strings.HasPrefix(m.name, `Test`) &&
		m.IsConcreteFunc() &&
		m.hasTestParam(`T`)
}

func (m *methodImp) IsBenchmark() bool {
	return strings.HasPrefix(m.name, `Benchmark`) &&
		m.IsConcreteFunc() &&
		m.hasTestParam(`B`)
}

func (m *methodImp) IsFuzz() bool {
	return strings.HasPrefix(m.name, `Fuzz`) &&
		m.IsConcreteFunc() &&
		m.hasTestParam(`F`)
}

func (m *methodImp) IsExample() bool {
	return strings.HasPrefix(m.name, `Example`) &&
		m.IsConcreteFunc() &&
		m.signature.IsVacant()
}

func (m *methodImp) IsConcreteFunc() bool {
	return !m.HasReceiver() &&
		!m.IsGeneric()
}

func (m *methodImp) IsNamed() bool {
	return len(m.name) > 0
}

func (m *methodImp) IsGeneric() bool {
	return len(m.typeParams) > 0 || (!utils.IsNil(m.receiver) && m.receiver.IsGeneric())
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
		return jsonify.New(ctx, m.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, m.Kind(), m.Index())
	}
	if ctx.SkipDead() && !m.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && m.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, m.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, m.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, m.Alive()).
		Add(ctx.OnlyIndex(), `package`, m.pkg).
		Add(ctx, `name`, m.name).
		AddNonZero(ctx, `loc`, m.loc).
		AddNonZeroIf(ctx, m.exported, `vis`, `exported`).
		AddNonZero(ctx.OnlyIndex(), `typeParams`, m.typeParams).
		Add(ctx.OnlyIndex(), `signature`, m.signature).
		AddNonZero(ctx.OnlyIndex(), `metrics`, m.metrics).
		AddNonZero(ctx.OnlyIndex(), `instances`, constructs.JsonSet(ctx.OnlyIndex(), m.instances.ToSlice())).
		AddNonZero(ctx.OnlyIndex(), `receiver`, m.receiver).
		AddNonZero(ctx, `ptrRecv`, m.ptrRecv).
		AddNonZeroIf(ctx, ctx.IsDebugReceiverIncluded(), `recvName`, m.recvName)
}

func (m *methodImp) ToStringer(s stringer.Stringer) {
	s.Write(`func `, m.pkg.Path(), `.`)
	if len(m.recvName) > 0 {
		s.Write(m.recvName, `.`)
	}
	s.Write(m.name).WriteList(`[`, `, `, `]`, m.typeParams)

	pre := s.String()
	s.Write(m.signature)
	suffix := s.String()[len(pre):]
	suffix, _ = strings.CutPrefix(suffix, `func`)
	s.Reset().Write(pre, suffix)
}

func (m *methodImp) String() string {
	return stringer.String(m)
}
