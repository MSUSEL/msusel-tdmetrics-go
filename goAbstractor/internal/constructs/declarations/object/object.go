package object

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/components/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

const Kind = `object`

// Object is a named type typically explicitly defined at the given location
// in the source code. An object typically handles structs with optional
// parameter types. An object can handle any type that methods can use
// as a receiver.
//
// If type parameters are given then the object is generic.
// Instances with realized versions of the object,
// are added for each used instance in the source code.
// If there are no instances then the generic object isn't used.
type Object interface {
	declarations.Declaration
	_object()

	AddMethod(met method.Method) method.Method
	AddInstance(inst instance.Instance) instance.Instance
	SetInterface(it interfaceDesc.InterfaceDesc)

	IsNamed() bool
	IsGeneric() bool
}

type Args struct {
	RealType types.Type
	Package  declarations.Package
	Name     string
	Location locs.Loc

	TypeParams []typeParam.TypeParam
	Data       structDesc.StructDesc
}

type objectImp struct {
	realType types.Type
	pkg      declarations.Package
	name     string
	loc      locs.Loc

	typeParams []typeParam.TypeParam
	data       structDesc.StructDesc

	methods   collections.SortedSet[method.Method]
	instances collections.SortedSet[instance.Instance]
	inter     interfaceDesc.InterfaceDesc

	index int
}

func newObject(args Args) Object {
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

		methods:   sortedSet.New[method.Method](),
		instances: sortedSet.New[instance.Instance](),
	}
}

func (d *objectImp) _object()           {}
func (d *objectImp) Kind() string       { return Kind }
func (d *objectImp) SetIndex(index int) { d.index = index }
func (d *objectImp) GoType() types.Type { return d.realType }

func (d *objectImp) Package() declarations.Package { return d.pkg }
func (d *objectImp) Name() string                  { return d.name }
func (d *objectImp) Location() locs.Loc            { return d.loc }

func (d *objectImp) AddMethod(met method.Method) method.Method {
	v, _ := d.methods.TryAdd(met)
	return v
}

func (d *objectImp) AddInstance(inst instance.Instance) instance.Instance {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *objectImp) SetInterface(it interfaceDesc.InterfaceDesc) {
	d.inter = it
}

func (d *objectImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *objectImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *objectImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[Object](d, other, Comparer())
}

func Comparer() comp.Comparer[Object] {
	return func(a, b Object) int {
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
