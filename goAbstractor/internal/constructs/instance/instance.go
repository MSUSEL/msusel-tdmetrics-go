package instance

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type instanceImp struct {
	realType      types.Type
	generic       constructs.Declaration
	resolved      constructs.TypeDesc
	instanceTypes []constructs.TypeDesc

	index int
}

func newInstance(args constructs.InstanceArgs) constructs.Instance {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`instance types`, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	if utils.IsNil(args.RealType) {
		pkg := args.Generic.Package()
		assert.ArgNotNil(`package`, pkg)

		// TODO: Implement if needed.
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

func (i *instanceImp) IsInstance() {}
func (i *instanceImp) IsTypeDesc() {}

func (i *instanceImp) Kind() kind.Kind    { return kind.Instance }
func (i *instanceImp) SetIndex(index int) { i.index = index }
func (m *instanceImp) GoType() types.Type { return m.realType }

func (m *instanceImp) Generic() constructs.Declaration { return m.generic }
func (m *instanceImp) Resolved() constructs.TypeDesc   { return m.resolved }

func (m *instanceImp) InstanceTypes() []constructs.TypeDesc {
	return m.instanceTypes
}

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Instance](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Instance] {
	return func(a, b constructs.Instance) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		return comp.Or(
			constructs.ComparerPend(aImp.resolved, bImp.resolved),
			constructs.SliceComparerPend(bImp.instanceTypes, bImp.instanceTypes),
		)
	}
}

func (i *instanceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, i.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, i.index).
		AddNonZero(ctx2, `generic`, i.generic).
		AddNonZero(ctx2, `resolved`, i.resolved).
		AddNonZero(ctx2, `instanceTypes`, i.instanceTypes)
}
