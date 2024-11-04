package interfaceInst

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type instanceImp struct {
	realType      types.Type
	generic       constructs.InterfaceDecl
	resolved      constructs.InterfaceDesc
	instanceTypes []constructs.TypeDesc
	index         int
	alive         bool
}

func newInstance(args constructs.InterfaceInstArgs) constructs.InterfaceInst {
	assert.ArgNotNil(`generic`, args.Generic)
	assert.ArgNotNil(`resolved`, args.Resolved)
	assert.ArgNotEmpty(`instance types`, args.InstanceTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)
	if !args.Generic.IsGeneric() {
		panic(terror.New(`may not create an instance on a non-generic interface`).
			With(`interface`, args.Generic))
	}

	if utils.IsNil(args.RealType) {
		pkg := args.Generic.Package()
		assert.ArgNotNil(`package`, pkg)

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
		instanceTypes: args.InstanceTypes,
	}
	return args.Generic.AddInstance(inst)
}

func (i *instanceImp) IsInterfaceInst() {}
func (i *instanceImp) IsTypeDesc()      {}

func (i *instanceImp) Kind() kind.Kind     { return kind.InterfaceInst }
func (i *instanceImp) Index() int          { return i.index }
func (i *instanceImp) SetIndex(index int)  { i.index = index }
func (i *instanceImp) Alive() bool         { return i.alive }
func (i *instanceImp) SetAlive(alive bool) { i.alive = alive }
func (m *instanceImp) GoType() types.Type  { return m.realType }

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
			constructs.SliceComparerPend(aImp.instanceTypes, bImp.instanceTypes),
			constructs.ComparerPend(aImp.generic, bImp.generic),
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
		Add(ctx.OnlyIndex(), `generic`, i.generic).
		Add(ctx.OnlyIndex(), `resolved`, i.resolved).
		Add(ctx.Short(), `instanceTypes`, i.instanceTypes)
}

func (i *instanceImp) String() string {
	return i.generic.Name() +
		`[` + enumerator.Enumerate(i.instanceTypes...).Join(`, `) + `]` +
		i.resolved.String()
}
