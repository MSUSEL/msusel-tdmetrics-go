package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

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
	Declaration
	TypeDesc
	_object()

	Package() Package
	Name() string
	Location() locs.Loc

	addMethod(met Method) Method
	addInstance(inst Instance) Instance
	addImplements(it Interface) Interface

	IsNamed() bool
	IsGeneric() bool
}

type ObjectArgs struct {
	RealType types.Type
	Package  Package
	Name     string
	Location locs.Loc

	TypeParams []TypeParam
	Fields     []Field
}

type objectImp struct {
	realType types.Type
	pkg      Package
	name     string
	loc      locs.Loc

	typeParams []TypeParam
	fields     []Field

	instances  Set[Instance]
	methods    Set[Method]
	implements Set[Interface]

	index int
}

func newObject(args ObjectArgs) Object {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNoNils(`fields`, args.Fields)
	assert.ArgNoNils(`type params`, args.TypeParams)

	return &objectImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		typeParams: args.TypeParams,
		fields:     args.Fields,

		instances:  NewSet[Instance](),
		methods:    NewSet[Method](),
		implements: NewSet[Interface](),
	}
}

func (d *objectImp) _object()           {}
func (d *objectImp) Kind() kind.Kind    { return kind.Object }
func (d *objectImp) setIndex(index int) { d.index = index }
func (d *objectImp) GoType() types.Type { return d.realType }

func (d *objectImp) Package() Package   { return d.pkg }
func (d *objectImp) Name() string       { return d.name }
func (d *objectImp) Location() locs.Loc { return d.loc }

func (d *objectImp) addMethod(met Method) Method {
	return d.methods.Insert(met)
}

func (d *objectImp) addInstance(inst Instance) Instance {
	return d.instances.Insert(inst)
}

func (d *objectImp) addImplements(it Interface) Interface {
	return d.implements.Insert(it)
}

func addInheritance(roots []Object, decl Object) []Object {

	print(decl)
	// TODO: Implement

	return roots
}

func (d *objectImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *objectImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *objectImp) compareTo(other Construct) int {
	b := other.(*objectImp)
	return or(
		func() int { return Compare(d.pkg, b.pkg) },
		func() int { return strings.Compare(d.name, b.name) },
		func() int { return compareSlice(d.typeParams, b.typeParams) },
		func() int { return compareSlice(d.fields, b.fields) },
	)
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
		AddNonZero(ctx2, `fields`, d.fields).
		AddNonZero(ctx2, `instances`, d.instances).
		AddNonZero(ctx2, `methods`, d.methods).
		AddNonZero(ctx2, `implements`, d.implements)
}
