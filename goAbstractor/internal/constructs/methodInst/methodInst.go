package methodInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type instanceImp struct {
	realType      types.Type
	generic       constructs.Method
	resolved      constructs.Signature
	instanceTypes []constructs.TypeDesc
	receiver      constructs.ObjectInst
	id            any
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

func (i *instanceImp) Kind() kind.Kind { return kind.MethodInst }
func (i *instanceImp) Id() any         { return i.id }
func (i *instanceImp) SetId(id any)    { i.id = id }

func (m *instanceImp) Generic() constructs.Method     { return m.generic }
func (m *instanceImp) Resolved() constructs.Signature { return m.resolved }

func (m *instanceImp) InstanceTypes() []constructs.TypeDesc  { return m.instanceTypes }
func (m *instanceImp) Receiver() constructs.ObjectInst       { return m.receiver }
func (m *instanceImp) SetReceiver(obj constructs.ObjectInst) { m.receiver = obj }

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.MethodInst](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.MethodInst] {
	return func(a, b constructs.MethodInst) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		return comp.Or(
			constructs.ComparerPend(aImp.resolved, bImp.resolved),
			constructs.SliceComparerPend(bImp.instanceTypes, bImp.instanceTypes),
		)
	}
}

func (i *instanceImp) RemoveTempReferences() {
	for j, it := range i.instanceTypes {
		i.instanceTypes[j] = constructs.ResolvedTempReference(it)
	}
}

func (i *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, i.id)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsIdShown(), `id`, i.id).
		AddNonZero(ctx2, `generic`, i.generic).
		AddNonZero(ctx2, `resolved`, i.resolved).
		AddNonZero(ctx2, `instanceTypes`, i.instanceTypes).
		AddNonZero(ctx2, `receiver`, i.receiver)
}

func (i *instanceImp) String() string {
	return i.generic.Name() +
		`[` + enumerator.Enumerate(i.instanceTypes...).Join(`, `) + `]` +
		i.resolved.String()
}
