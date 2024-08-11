package object

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

const Kind = `object`

type Args struct {
	RealType types.Type
	Package  constructs.Package
	Name     string
	Location locs.Loc

	TypeParams []typeDescs.TypeParam
	Data       typeDescs.StructDesc
}

type objectImp struct {
	realType types.Type
	pkg      constructs.Package
	name     string
	loc      locs.Loc

	typeParams []typeDescs.TypeParam
	data       typeDescs.StructDesc

	methods   collections.SortedSet[declarations.Method]
	instances collections.SortedSet[components.Instance]
	inter     typeDescs.InterfaceDesc

	index int
}

func New(args Args) declarations.Object {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgNoNils(`type params`, args.TypeParams)
	assert.ArgNotNil(`data`, args.Data)

	return &objectImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		typeParams: args.TypeParams,
		data:       args.Data,

		methods:   sortedSet.New[declarations.Method](),
		instances: sortedSet.New[components.Instance](),
	}
}

func (d *objectImp) IsDeclaration()     {}
func (d *objectImp) IsTypeDesc()        {}
func (d *objectImp) IsObject()          {}
func (d *objectImp) Kind() string       { return Kind }
func (d *objectImp) SetIndex(index int) { d.index = index }
func (d *objectImp) GoType() types.Type { return d.realType }

func (d *objectImp) Package() constructs.Package { return d.pkg }
func (d *objectImp) Name() string                { return d.name }
func (d *objectImp) Location() locs.Loc          { return d.loc }

func (d *objectImp) AddMethod(met declarations.Method) declarations.Method {
	v, _ := d.methods.TryAdd(met)
	return v
}

func (d *objectImp) AddInstance(inst components.Instance) components.Instance {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *objectImp) SetInterface(it typeDescs.InterfaceDesc) {
	d.inter = it
}

func (d *objectImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *objectImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *objectImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[declarations.Object](d, other, Comparer())
}

func Comparer() comp.Comparer[declarations.Object] {
	return func(a, b declarations.Object) int {
		aImp, bImp := a.(*objectImp), b.(*objectImp)
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.typeParams, bImp.typeParams),
			constructs.ComparerPend(aImp.data, bImp.data),
		)
	}
}

func (d *objectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, d.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, d.index).
		AddNonZero(ctx2, `package`, d.pkg).
		AddNonZero(ctx2, `name`, d.name).
		AddNonZero(ctx2, `loc`, d.loc).
		AddNonZero(ctx2, `typeParams`, d.typeParams).
		AddNonZero(ctx2, `data`, d.data).
		AddNonZero(ctx2, `instances`, d.instances).
		AddNonZero(ctx2, `methods`, d.methods).
		AddNonZero(ctx2, `interface`, d.inter)
}
