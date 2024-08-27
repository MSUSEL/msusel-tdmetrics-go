package objectInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type instanceImp struct {
	realType      types.Type
	generic       constructs.Object
	resolved      constructs.StructDesc
	instanceTypes []constructs.TypeDesc
	methods       collections.SortedSet[constructs.MethodInst]
	inter         constructs.InterfaceDesc
	id            any
}

func newInstance(args constructs.ObjectInstArgs) constructs.ObjectInst {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`instance types`, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	if utils.IsNil(args.RealType) {
		pkg := args.Generic.Package()
		assert.ArgNotNil(`package`, pkg)

		// Implement if needed.
		assert.NotImplemented()
	}
	assert.ArgNotNil(`real type`, args.RealType)

	inst := &instanceImp{
		realType:      args.RealType,
		generic:       args.Generic,
		resolved:      args.Resolved,
		instanceTypes: args.InstanceTypes,
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsObjectInst() {}
func (i *instanceImp) IsTypeDesc()   {}

func (i *instanceImp) Kind() kind.Kind    { return kind.ObjectInst }
func (i *instanceImp) Id() any            { return i.id }
func (i *instanceImp) SetId(id any)       { i.id = id }
func (m *instanceImp) GoType() types.Type { return m.realType }

func (m *instanceImp) Generic() constructs.Object           { return m.generic }
func (m *instanceImp) Resolved() constructs.StructDesc      { return m.resolved }
func (m *instanceImp) InstanceTypes() []constructs.TypeDesc { return m.instanceTypes }

func (m *instanceImp) AddMethods(method constructs.MethodInst) constructs.MethodInst {
	v, _ := m.methods.TryAdd(method)
	return v
}

func (m *instanceImp) Methods() collections.ReadonlySortedSet[constructs.MethodInst] {
	return m.methods.Readonly()
}

func (m *instanceImp) Interface() constructs.InterfaceDesc      { return m.inter }
func (m *instanceImp) SetInterface(it constructs.InterfaceDesc) { m.inter = it }

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.ObjectInst](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.ObjectInst] {
	return func(a, b constructs.ObjectInst) int {
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
		AddNonZero(ctx2, `methods`, i.methods.ToList())
}

func (i *instanceImp) String() string {
	return i.generic.Name() +
		`[` + enumerator.Enumerate(i.instanceTypes...).Join(`, `) + `]` +
		i.resolved.String()
}
