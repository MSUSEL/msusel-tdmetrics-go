package interfaceInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type instanceImp struct {
	constructs.ConstructCore
	realType      types.Type
	generic       constructs.InterfaceDecl
	resolved      constructs.InterfaceDesc
	implicitTypes []constructs.TypeDesc
	instanceTypes []constructs.TypeDesc
}

func newInstance(args constructs.InterfaceInstArgs) constructs.InterfaceInst {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.AnyArgNotEmpty(`implicit & instance types`, args.ImplicitTypes, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)
	if !args.Generic.IsGeneric() {
		panic(terror.New(`may not create an instance on a non-generic interface`).
			With(`interface`, args.Generic))
	}

	if utils.IsNil(args.RealType) {
		pkg := args.Generic.Package()
		assert.ArgNotNil(`package`, pkg)

		// TODO: Implement this for nested types.
		assert.ArgNotEmpty(`instance types`, args.InstanceTypes)
		assert.ArgIsEmpty(`implicit types`, args.ImplicitTypes)

		tArgs := make([]types.Type, len(args.InstanceTypes))
		for i, ip := range args.InstanceTypes {
			tArgs[i] = ip.GoType()
		}

		switch args.Generic.Interface().Hint() {
		case hint.Pointer:
			if len(tArgs) != 1 {
				panic(terror.New(`an instance of a pointer must have exactly one type argument`).
					With(`type args`, tArgs))
			}
			args.RealType = types.NewPointer(tArgs[0])

		case hint.List:
			if len(tArgs) != 1 {
				panic(terror.New(`an instance of a list must have exactly one type argument`).
					With(`type args`, tArgs))
			}
			args.RealType = types.NewSlice(tArgs[0])

		case hint.Map:
			if len(tArgs) != 2 {
				panic(terror.New(`an instance of a map must have exactly two type arguments`).
					With(`type args`, tArgs))
			}
			args.RealType = types.NewMap(tArgs[0], tArgs[1])

		case hint.Chan:
			if len(tArgs) != 1 {
				panic(terror.New(`an instance of a channel must have exactly one type argument`).
					With(`type args`, tArgs))
			}
			args.RealType = types.NewChan(types.SendRecv, tArgs[0])

		default:
			gt := args.Generic.GoType()
			ggt, ok := gt.(interface {
				types.Type
				TypeParams() *types.TypeParamList
			})
			if !ok {
				panic(terror.New(`go type is not a generic type`).
					With(`type`, args.Generic).
					With(`goType`, gt))
			}
			if ggt.TypeParams().Len() <= 0 {
				panic(terror.New(`may not create an instance with a non-generic go type`).
					With(`type`, args.Generic).
					With(`goType`, gt))
			}

			rt, err := types.Instantiate(nil, ggt, tArgs, true)
			if err != nil {
				panic(terror.New(`failed to instantiate an interface instance`, err))
			}
			args.RealType = rt
		}
	}
	assert.ArgNotNil(`real type`, args.RealType)

	inst := &instanceImp{
		realType:      args.RealType,
		generic:       args.Generic,
		resolved:      args.Resolved,
		implicitTypes: args.ImplicitTypes,
		instanceTypes: args.InstanceTypes,
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsInterfaceInst() {}
func (i *instanceImp) IsTypeDesc()      {}

func (i *instanceImp) Kind() kind.Kind    { return kind.InterfaceInst }
func (m *instanceImp) GoType() types.Type { return m.realType }

func (m *instanceImp) Generic() constructs.InterfaceDecl    { return m.generic }
func (m *instanceImp) Resolved() constructs.InterfaceDesc   { return m.resolved }
func (m *instanceImp) ImplicitTypes() []constructs.TypeDesc { return m.implicitTypes }
func (m *instanceImp) InstanceTypes() []constructs.TypeDesc { return m.instanceTypes }

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
	if ctx.SkipDuplicates() && i.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, i.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, i.Index()).
		Add(ctx.OnlyIndex(), `generic`, i.generic).
		Add(ctx.OnlyIndex(), `resolved`, i.resolved).
		AddNonZero(ctx.Short(), `implicitTypes`, i.implicitTypes).
		AddNonZero(ctx.Short(), `instanceTypes`, i.instanceTypes)
}

func (i *instanceImp) ToStringer(s stringer.Stringer) {
	s.Write(i.generic.Name(), `[`).
		WriteList(``, `, `, `;`, i.implicitTypes).
		WriteList(``, `, `, ``, i.instanceTypes).
		Write(`]`, i.resolved)
}

func (i *instanceImp) String() string {
	return stringer.String(i)
}
