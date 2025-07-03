package methodInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type instanceImp struct {
	constructs.ConstructCore
	generic       constructs.Method
	resolved      constructs.Signature
	instanceTypes []constructs.TypeDesc
	metrics       constructs.Metrics
	receiver      constructs.ObjectInst
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
		metrics:       args.Metrics,
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsMethodInst() {}

func (i *instanceImp) Kind() kind.Kind { return kind.MethodInst }
func (i *instanceImp) Name() string    { return i.generic.Name() }

func (i *instanceImp) FuncType() *types.Func                 { return i.generic.FuncType() }
func (i *instanceImp) SetMetrics(metrics constructs.Metrics) { i.metrics = metrics }
func (i *instanceImp) Metrics() constructs.Metrics           { return i.metrics }
func (i *instanceImp) Generic() constructs.Method            { return i.generic }
func (i *instanceImp) Resolved() constructs.Signature        { return i.resolved }
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

func (i *instanceImp) RemoveTempReferences(required bool) bool {
	changed := false
	var subChanged bool
	for j, it := range i.instanceTypes {
		i.instanceTypes[j], subChanged = constructs.ResolvedTempReference(it, required)
		changed = changed || subChanged
	}
	return changed
}

func (i *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, i.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, i.Kind(), i.Index())
	}
	if ctx.SkipDuplicates() && i.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, i.Index()).
		Add(ctx.OnlyIndex(), `generic`, i.generic).
		Add(ctx.OnlyIndex(), `resolved`, i.resolved).
		Add(ctx.Short(), `instanceTypes`, i.instanceTypes).
		AddNonZero(ctx.OnlyIndex(), `metrics`, i.metrics).
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
