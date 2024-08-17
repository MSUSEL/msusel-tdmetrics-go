package object

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type objectImp struct {
	realType types.Type
	pkg      constructs.Package
	name     string
	loc      locs.Loc

	typeParams []constructs.TypeParam
	data       constructs.StructDesc
	inter      constructs.InterfaceDesc

	methods   collections.SortedSet[constructs.Method]
	instances collections.SortedSet[constructs.Instance]

	index int
}

func newObject(args constructs.ObjectArgs) constructs.Object {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`type params`, args.TypeParams)
	assert.ArgNotNil(`data`, args.Data)

	return &objectImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		typeParams: args.TypeParams,
		data:       args.Data,

		methods:   sortedSet.New(method.Comparer()),
		instances: sortedSet.New(instance.Comparer()),
	}
}

func (d *objectImp) IsDeclaration()     {}
func (d *objectImp) IsTypeDesc()        {}
func (d *objectImp) IsObject()          {}
func (d *objectImp) Kind() kind.Kind    { return kind.Object }
func (d *objectImp) SetIndex(index int) { d.index = index }
func (d *objectImp) GoType() types.Type { return d.realType }

func (d *objectImp) Package() constructs.Package { return d.pkg }
func (d *objectImp) Name() string                { return d.name }
func (d *objectImp) Location() locs.Loc          { return d.loc }

func (d *objectImp) TypeParams() []constructs.TypeParam {
	return d.typeParams
}

func (d *objectImp) AddMethod(met constructs.Method) constructs.Method {
	v, _ := d.methods.TryAdd(met)
	return v
}

func (d *objectImp) AddInstance(inst constructs.Instance) constructs.Instance {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *objectImp) SetInterface(it constructs.InterfaceDesc) {
	d.inter = it
}

func (d *objectImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *objectImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *objectImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Object](d, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Object] {
	return func(a, b constructs.Object) int {
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
		AddIf(ctx, ctx.IsKindShown(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, d.index).
		AddNonZero(ctx2, `package`, d.pkg).
		AddNonZero(ctx2, `name`, d.name).
		AddNonZero(ctx2, `loc`, d.loc).
		AddNonZero(ctx2, `typeParams`, d.typeParams).
		AddNonZero(ctx2, `data`, d.data).
		AddNonZero(ctx2, `instances`, d.instances.ToSlice()).
		AddNonZero(ctx2, `methods`, d.methods.ToSlice()).
		AddNonZero(ctx2, `interface`, d.inter)
}
