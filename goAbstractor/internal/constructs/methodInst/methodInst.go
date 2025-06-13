package methodInst

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type instanceImp struct {
	generic       constructs.Method
	resolved      constructs.Signature
	instanceTypes []constructs.TypeDesc
	receiver      constructs.ObjectInst
	index         int
	alive         bool
}

func newInstance(args constructs.MethodInstArgs) constructs.MethodInst {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`instance types`, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	inst := &instanceImp{
		generic:       args.Generic,
		resolved:      args.Resolved,
		instanceTypes: args.InstanceTypes,
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsMethodInst() {}

func (i *instanceImp) Kind() kind.Kind     { return kind.MethodInst }
func (i *instanceImp) Index() int          { return i.index }
func (i *instanceImp) SetIndex(index int)  { i.index = index }
func (i *instanceImp) Alive() bool         { return i.alive }
func (i *instanceImp) SetAlive(alive bool) { i.alive = alive }
func (i *instanceImp) Name() string        { return i.generic.Name() }

func (i *instanceImp) Generic() constructs.Method     { return i.generic }
func (i *instanceImp) Resolved() constructs.Signature { return i.resolved }

func (i *instanceImp) TypeParams() []constructs.TypeParam    { return i.generic.TypeParams() }
func (i *instanceImp) InstanceTypes() []constructs.TypeDesc  { return i.instanceTypes }
func (i *instanceImp) Receiver() constructs.ObjectInst       { return i.receiver }
func (i *instanceImp) SetReceiver(obj constructs.ObjectInst) { i.receiver = obj }

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.MethodInst](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.MethodInst] {
	return func(a, b constructs.MethodInst) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.resolved, bImp.resolved),
			constructs.SliceComparerPend(aImp.instanceTypes, bImp.instanceTypes),
			constructs.ComparerPend(aImp.generic, bImp.generic),
		)
	}
}

func (i *instanceImp) RemoveTempReferences(required bool) {
	for j, it := range i.instanceTypes {
		i.instanceTypes[j] = constructs.ResolvedTempReference(it, required)
	}
}

func (i *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, i.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, i.Kind(), i.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, i.index).
		Add(ctx.OnlyIndex(), `generic`, i.generic).
		Add(ctx.OnlyIndex(), `resolved`, i.resolved).
		Add(ctx.Short(), `instanceTypes`, i.instanceTypes).
		AddNonZero(ctx.OnlyIndex(), `receiver`, i.receiver)
}

func (i *instanceImp) ToStringer(s stringer.Stringer) {
	s.Write(i.generic.Name(), `[`).
		WriteList(``, `, `, ``, i.instanceTypes).
		Write(`]`, i.resolved)
}

func (i *instanceImp) String() string {
	return stringer.String(i)
}
