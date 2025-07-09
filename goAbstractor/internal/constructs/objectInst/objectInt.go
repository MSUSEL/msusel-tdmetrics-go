package objectInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/methodInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type instanceImp struct {
	constructs.ConstructCore
	realType          types.Type
	generic           constructs.Object
	resolvedData      constructs.StructDesc
	resolvedInterface constructs.InterfaceDesc
	implicitTypes     []constructs.TypeDesc
	instanceTypes     []constructs.TypeDesc
	methods           collections.SortedSet[constructs.MethodInst]
}

func newInstance(args constructs.ObjectInstArgs) constructs.ObjectInst {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved data`, args.ResolvedData)
	assert.AnyArgNotEmpty(`implicit & instance types`, args.ImplicitTypes, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	if utils.IsNil(args.RealType) {
		pkg := args.Generic.Package()
		assert.ArgNotNil(`package`, pkg)

		// Implement if needed.
		panic(terror.New(`not implemented`))
	}
	assert.ArgNotNil(`real type`, args.RealType)

	inst := &instanceImp{
		realType:      args.RealType,
		generic:       args.Generic,
		resolvedData:  args.ResolvedData,
		implicitTypes: args.ImplicitTypes,
		instanceTypes: args.InstanceTypes,
		methods:       sortedSet.New(methodInst.Comparer()),
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsObjectInst() {}
func (i *instanceImp) IsTypeDesc()   {}

func (i *instanceImp) Kind() kind.Kind    { return kind.ObjectInst }
func (m *instanceImp) GoType() types.Type { return m.realType }

func (m *instanceImp) Generic() constructs.Object                  { return m.generic }
func (m *instanceImp) ResolvedData() constructs.StructDesc         { return m.resolvedData }
func (m *instanceImp) ResolvedInterface() constructs.InterfaceDesc { return m.resolvedInterface }
func (m *instanceImp) ImplicitTypes() []constructs.TypeDesc        { return m.implicitTypes }
func (m *instanceImp) InstanceTypes() []constructs.TypeDesc        { return m.instanceTypes }

func (m *instanceImp) AddMethod(method constructs.MethodInst) constructs.MethodInst {
	v, _ := m.methods.TryAdd(method)
	return v
}

func (m *instanceImp) Methods() collections.ReadonlySortedSet[constructs.MethodInst] {
	return m.methods.Readonly()
}

func (m *instanceImp) SetResolvedInterface(it constructs.InterfaceDesc) {
	m.resolvedInterface = it
}

func (i *instanceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.ObjectInst](i, other, Comparer())
}

func Comparer() comp.Comparer[constructs.ObjectInst] {
	return func(a, b constructs.ObjectInst) int {
		aImp, bImp := a.(*instanceImp), b.(*instanceImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.resolvedData, bImp.resolvedData),
			constructs.SliceComparerPend(aImp.implicitTypes, bImp.implicitTypes),
			constructs.SliceComparerPend(aImp.instanceTypes, bImp.instanceTypes),
			constructs.ComparerPend(aImp.generic, bImp.generic),
		)
	}
}

func (i *instanceImp) RemoveTempReferences(required bool) bool {
	changed := false
	var subChanged bool
	for j, it := range i.implicitTypes {
		i.implicitTypes[j], subChanged = constructs.ResolvedTempReference(it, required)
		changed = changed || subChanged
	}
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
	if ctx.SkipDead() && !i.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && i.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, i.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, i.Alive()).
		Add(ctx.OnlyIndex(), `generic`, i.generic).
		Add(ctx.OnlyIndex(), `resData`, i.resolvedData).
		Add(ctx.OnlyIndex(), `resInterface`, i.resolvedInterface).
		AddNonZero(ctx.Short(), `implicitTypes`, i.implicitTypes).
		AddNonZero(ctx.Short(), `instanceTypes`, i.instanceTypes).
		AddNonZero(ctx.OnlyIndex(), `methods`, constructs.JsonSet(ctx.OnlyIndex(), i.methods.ToSlice()))
}

func (i *instanceImp) ToStringer(s stringer.Stringer) {
	s.Write(i.generic.Name(), `[`).
		WriteList(``, `, `, `;`, i.implicitTypes).
		WriteList(``, `, `, ``, i.instanceTypes).
		Write(`]`, i.resolvedData, i.resolvedInterface)
}

func (i *instanceImp) String() string {
	return stringer.String(i)
}
