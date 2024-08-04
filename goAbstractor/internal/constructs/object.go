package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

// Object is a named type explicitly defined at the given
// location in the source code. The underlying type description
// can be a class or interface with optional parameter types.
//
// If type parameters are given then the object is generic.
// Instances with realized versions of the class or interface,
// are added for each used instance in the source code. If there
// are no instances then this generic object isn't used.
//
// If the type description is an interface, then no methods will
// be added. Any method added to a class indicates that the
// class is the receiver for that method.
type Object interface {
	Declaration
	TypeDesc
	_object()

	Package() Package
	Name() string
	Location() locs.Loc

	addMethod(met Method) Method
	addInstance(inst Instance) Instance

	IsNamed() bool
	IsGeneric() bool
	IsInterface() bool
}

type ObjectArgs struct {
	RealType types.Type
	Package  Package
	Name     string
	Location locs.Loc

	Fields     []Field
	TypeParams []TypeParam

	// Exact types are like `string|int|bool` where the
	// data type must match exactly.
	Exact []TypeDesc

	// Approx types are like `~string|~int` where the data type
	// may be exact or an extension of the base type.
	Approx []TypeDesc
}

type objectImp struct {
	realType types.Type
	pkg      Package
	name     string
	loc      locs.Loc

	fields     []Field
	typeParams []TypeParam
	exact      []TypeDesc
	approx     []TypeDesc

	instances Set[Instance]
	methods   Set[Method]

	index int
}

func newObject(args ObjectArgs) Object {
	assert.ArgNotNil(`realType`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNoNils(`fields`, args.Fields)
	assert.ArgNoNils(`type params`, args.TypeParams)

	if len(args.Fields) > 0 && (len(args.Exact) > 0 || len(args.Approx) > 0) {
		panic(terror.New(`may not have fields and exact/approx interface values`).
			With(`fields`, len(args.Fields)).
			With(`exact`, len(args.Exact)).
			With(`approx`, len(args.Approx)))
	}

	return &objectImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		fields:     args.Fields,
		typeParams: args.TypeParams,
		exact:      args.Exact,
		approx:     args.Approx,

		instances: NewSet[Instance](),
		methods:   NewSet[Method](),
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

func (d *objectImp) IsInterface() bool {
	return len(d.fields) <= 0
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
		AddNonZero(ctx2, `fields`, d.fields).
		AddNonZero(ctx2, `typeParams`, d.typeParams).
		AddNonZero(ctx2, `approx`, d.approx).
		AddNonZero(ctx2, `exact`, d.exact).
		AddNonZero(ctx2, `instances`, d.instances).
		AddNonZero(ctx2, `methods`, d.methods)
}
