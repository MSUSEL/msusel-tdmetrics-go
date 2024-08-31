package interfaceInst

import (
	"go/types"

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
	generic       constructs.InterfaceDecl
	resolved      constructs.InterfaceDesc
	instanceTypes []constructs.TypeDesc
	index         int
}

func newInstance(args constructs.InterfaceInstArgs) constructs.InterfaceInst {
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

func (i *instanceImp) IsInterfaceInst() {}
func (i *instanceImp) IsTypeDesc()      {}

func (i *instanceImp) Kind() kind.Kind    { return kind.InterfaceInst }
func (i *instanceImp) Index() int         { return i.index }
func (i *instanceImp) SetIndex(index int) { i.index = index }
func (m *instanceImp) GoType() types.Type { return m.realType }

func (m *instanceImp) Generic() constructs.InterfaceDecl  { return m.generic }
func (m *instanceImp) Resolved() constructs.InterfaceDesc { return m.resolved }

func (m *instanceImp) InstanceTypes() []constructs.TypeDesc {
	return m.instanceTypes
}

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.InterfaceInst](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.InterfaceInst] {
	return func(a, b constructs.InterfaceInst) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		if aImp == bImp {
			return 0
		}
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
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, i.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, i.Kind(), i.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, i.index).
		AddNonZero(ctx.OnlyIndex(), `generic`, i.generic).
		AddNonZero(ctx.OnlyIndex(), `resolved`, i.resolved).
		AddNonZero(ctx.Short(), `instanceTypes`, i.instanceTypes)
}

func (i *instanceImp) String() string {
	return i.generic.Name() +
		`[` + enumerator.Enumerate(i.instanceTypes...).Join(`, `) + `]` +
		i.resolved.String()
}
