package interfaceDecl

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type interfaceDeclImp struct {
	realType types.Type
	pkg      constructs.Package
	name     string
	loc      locs.Loc

	typeParams []constructs.TypeParam
	inter      constructs.InterfaceDesc

	instances collections.SortedSet[constructs.Instance]

	index int
}

func newInterfaceDecl(args constructs.InterfaceDeclArgs) constructs.InterfaceDecl {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`interface`, args.Interface)
	assert.ArgHasNoNils(`type params`, args.TypeParams)

	if utils.IsNil(args.RealType) {
		pkg := args.Package.Source().Types
		assert.ArgNotNil(`package`, pkg)

		tn := types.NewTypeName(args.Location.Pos(), pkg, args.Name, nil)
		args.RealType = types.NewNamed(tn, args.Interface.GoType(), nil)
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &interfaceDeclImp{
		realType: args.RealType,
		pkg:      args.Package,
		name:     args.Name,
		loc:      args.Location,

		typeParams: args.TypeParams,
		inter:      args.Interface,

		instances: sortedSet.New(instance.Comparer()),
	}
}

func (d *interfaceDeclImp) IsDeclaration() {}
func (d *interfaceDeclImp) IsTypeDesc()    {}
func (d *interfaceDeclImp) IsInterface()   {}

func (d *interfaceDeclImp) Kind() kind.Kind    { return kind.InterfaceDecl }
func (d *interfaceDeclImp) SetIndex(index int) { d.index = index }
func (d *interfaceDeclImp) GoType() types.Type { return d.realType }

func (d *interfaceDeclImp) Package() constructs.Package { return d.pkg }
func (d *interfaceDeclImp) Name() string                { return d.name }
func (d *interfaceDeclImp) Location() locs.Loc          { return d.loc }

func (d *interfaceDeclImp) Type() constructs.TypeDesc           { return d.inter }
func (d *interfaceDeclImp) Interface() constructs.InterfaceDesc { return d.inter }
func (d *interfaceDeclImp) TypeParams() []constructs.TypeParam  { return d.typeParams }

func (d *interfaceDeclImp) Instances() collections.ReadonlySortedSet[constructs.Instance] {
	return d.instances.Readonly()
}

func (d *interfaceDeclImp) AddInstance(inst constructs.Instance) constructs.Instance {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *interfaceDeclImp) IsNamed() bool {
	return len(d.name) > 0
}

func (d *interfaceDeclImp) IsGeneric() bool {
	return len(d.typeParams) > 0
}

func (d *interfaceDeclImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.InterfaceDecl](d, other, Comparer())
}

func Comparer() comp.Comparer[constructs.InterfaceDecl] {
	return func(a, b constructs.InterfaceDecl) int {
		aImp, bImp := a.(*interfaceDeclImp), b.(*interfaceDeclImp)
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.typeParams, bImp.typeParams),
			constructs.ComparerPend(aImp.inter, bImp.inter),
		)
	}
}

func (d *interfaceDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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
		AddNonZero(ctx2, `interface`, d.inter).
		AddNonZero(ctx2, `instances`, d.instances.ToSlice())
}

func (d *interfaceDeclImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(d.name + ` `)
	if len(d.typeParams) > 0 {
		buf.WriteString(`[` + enumerator.Enumerate(d.typeParams).Join(`, `) + `]`)
	}
	buf.WriteString(d.inter.String())
	return buf.String()
}
