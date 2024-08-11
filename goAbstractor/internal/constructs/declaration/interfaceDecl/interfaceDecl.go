package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// Interface is a named interface typically explicitly defined at the given
// location in the source code. The underlying type description
// can be a class or interface with optional parameter types.
//
// If type parameters are given then the interface is generic.
// Instances with realized versions of the interface,
// are added for each used instance in the source code. If there
// are no instances then the generic interface isn't used.
type Interface interface {
	Declaration
	TypeDesc
	_interface()

	Package() Package
	Name() string
	Location() locs.Loc

	addInstance(inst Instance) Instance
	addInherits(it Interface) Interface

	IsNamed() bool
	IsGeneric() bool
}

type InterfaceArgs struct {
	RealType types.Type
	Package  Package
	Name     string
	Location locs.Loc

	TypeParams []TypeParam

	// Methods is the set of signatures for this interface.
	Methods []Method

	// Exact types are like `string|int|bool` where the
	// data type must match exactly.
	Exact []TypeDesc

	// Approx types are like `~string|~int` where the data type
	// may be exact or an extension of the base type.
	Approx []TypeDesc
}

type interfaceImp struct {
	realType types.Type
	pkg      Package
	name     string
	loc      locs.Loc

	typeParams []TypeParam
	methods    []Method
	exact      []TypeDesc
	approx     []TypeDesc

	instances Set[Instance]
	inherits  Set[Interface]

	index int
}

func newInterface(args InterfaceArgs) Interface {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNoNils(`methods`, args.Methods)
	assert.ArgNoNils(`type params`, args.TypeParams)

	// TODO: Check that none of the methods have receivers.

	return &interfaceImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		typeParams: args.TypeParams,
		methods:    args.Methods,
		exact:      args.Exact,
		approx:     args.Approx,

		instances: NewSet[Instance](),
		inherits:  NewSet[Interface](),
	}
}

func (d *interfaceImp) _interface()        {}
func (d *interfaceImp) Kind() kind.Kind    { return kind.Interface }
func (d *interfaceImp) setIndex(index int) { d.index = index }
func (d *interfaceImp) GoType() types.Type { return d.realType }

func (d *interfaceImp) Package() Package   { return d.pkg }
func (d *interfaceImp) Name() string       { return d.name }
func (d *interfaceImp) Location() locs.Loc { return d.loc }

func (d *interfaceImp) addInstance(inst Instance) Instance {
	return d.instances.Insert(inst)
}

func (d *interfaceImp) addInherits(it Interface) Interface {
	return d.inherits.Insert(it)
}

func (d *interfaceImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *interfaceImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *interfaceImp) compareTo(other Construct) int {
	b := other.(*interfaceImp)
	return or(
		func() int { return Compare(d.pkg, b.pkg) },
		func() int { return strings.Compare(d.name, b.name) },
		func() int { return compareSlice(d.typeParams, b.typeParams) },
		func() int { return compareSlice(d.methods, b.methods) },
		func() int { return compareSlice(d.exact, b.exact) },
		func() int { return compareSlice(d.approx, b.approx) },
	)
}

func (d *interfaceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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
		AddNonZero(ctx2, `methods`, d.methods).
		AddNonZero(ctx2, `approx`, d.approx).
		AddNonZero(ctx2, `exact`, d.exact).
		AddNonZero(ctx2, `instances`, d.instances).
		AddNonZero(ctx2, `inherits`, d.inherits)
}
