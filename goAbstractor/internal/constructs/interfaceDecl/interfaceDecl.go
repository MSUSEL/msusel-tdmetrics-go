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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type interfaceDeclImp struct {
	realType types.Type
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc
	index    int
	alive    bool

	typeParams []constructs.TypeParam
	inter      constructs.InterfaceDesc
	instances  collections.SortedSet[constructs.InterfaceInst]
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
		realType:   args.RealType,
		pkg:        args.Package,
		name:       args.Name,
		exported:   args.Exported,
		loc:        args.Location,
		typeParams: args.TypeParams,
		inter:      args.Interface,
		instances:  sortedSet.New(interfaceInst.Comparer()),
	}
}

func (d *interfaceDeclImp) IsDeclaration() {}
func (d *interfaceDeclImp) IsTypeDesc()    {}
func (d *interfaceDeclImp) IsInterface()   {}

func (d *interfaceDeclImp) Kind() kind.Kind     { return kind.InterfaceDecl }
func (d *interfaceDeclImp) Index() int          { return d.index }
func (d *interfaceDeclImp) SetIndex(index int)  { d.index = index }
func (d *interfaceDeclImp) Alive() bool         { return d.alive }
func (d *interfaceDeclImp) SetAlive(alive bool) { d.alive = alive }
func (d *interfaceDeclImp) GoType() types.Type  { return d.realType }
func (d *interfaceDeclImp) Name() string        { return d.name }
func (d *interfaceDeclImp) Exported() bool      { return d.exported }
func (d *interfaceDeclImp) Location() locs.Loc  { return d.loc }

func (d *interfaceDeclImp) Package() constructs.Package         { return d.pkg }
func (d *interfaceDeclImp) Type() constructs.TypeDesc           { return d.inter }
func (d *interfaceDeclImp) Interface() constructs.InterfaceDesc { return d.inter }
func (d *interfaceDeclImp) TypeParams() []constructs.TypeParam  { return d.typeParams }

func (d *interfaceDeclImp) Instances() collections.ReadonlySortedSet[constructs.InterfaceInst] {
	return d.instances.Readonly()
}

func (d *interfaceDeclImp) AddInstance(inst constructs.InterfaceInst) constructs.InterfaceInst {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *interfaceDeclImp) FindInstance(instanceTypes []constructs.TypeDesc) (constructs.InterfaceInst, bool) {
	cmp := constructs.SliceComparer[constructs.TypeDesc]()
	return d.instances.Enumerate().Where(func(i constructs.InterfaceInst) bool {
		return cmp(instanceTypes, i.InstanceTypes()) == 0
	}).First()
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
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.typeParams, bImp.typeParams),
			constructs.ComparerPend(aImp.inter, bImp.inter),
		)
	}
}

func (d *interfaceDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, d.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, d.Kind(), d.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, d.index).
		Add(ctx.OnlyIndex(), `package`, d.pkg).
		Add(ctx, `name`, d.name).
		Add(ctx.OnlyIndex(), `interface`, d.inter).
		AddNonZero(ctx, `loc`, d.loc).
		AddNonZeroIf(ctx, d.exported, `vis`, `exported`).
		AddNonZero(ctx.OnlyIndex(), `typeParams`, d.typeParams).
		AddNonZero(ctx.OnlyIndex(), `instances`, d.instances.ToSlice())
}

func (d *interfaceDeclImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(d.pkg.Path())
	buf.WriteString(`.`)
	buf.WriteString(d.name)
	if len(d.typeParams) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(d.typeParams...).Join(`, `))
		buf.WriteString(`]`)
	}
	buf.WriteString(` interface{--}`)
	return buf.String()
}
