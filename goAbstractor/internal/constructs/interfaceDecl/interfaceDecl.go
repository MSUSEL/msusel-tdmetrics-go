package interfaceDecl

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type interfaceDeclImp struct {
	constructs.ConstructCore
	realType types.Type
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc
	nest     constructs.NestType

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
		nest:       args.Nest,
		instances:  sortedSet.New(interfaceInst.Comparer()),
	}
}

func (d *interfaceDeclImp) IsDeclaration() {}
func (d *interfaceDeclImp) IsTypeDesc()    {}
func (d *interfaceDeclImp) IsInterface()   {}

func (d *interfaceDeclImp) Kind() kind.Kind    { return kind.InterfaceDecl }
func (d *interfaceDeclImp) GoType() types.Type { return d.realType }
func (d *interfaceDeclImp) Name() string       { return d.name }
func (d *interfaceDeclImp) Exported() bool     { return d.exported }
func (d *interfaceDeclImp) Location() locs.Loc { return d.loc }

func (d *interfaceDeclImp) Package() constructs.Package         { return d.pkg }
func (d *interfaceDeclImp) Type() constructs.TypeDesc           { return d.inter }
func (d *interfaceDeclImp) Interface() constructs.InterfaceDesc { return d.inter }
func (d *interfaceDeclImp) TypeParams() []constructs.TypeParam  { return d.typeParams }
func (d *interfaceDeclImp) Nest() constructs.NestType           { return d.nest }

func (d *interfaceDeclImp) ImplicitTypeParams() []constructs.TypeParam {
	if d.nest == nil {
		return nil
	}
	if method, ok := d.nest.(constructs.Method); ok {
		return method.TypeParams()
	}
	panic(terror.New(`may not get ImplicitTypeParams from a non-method nesting declaration`).
		WithType(`nest`, d.nest))
}

func (d *interfaceDeclImp) Instances() collections.ReadonlySortedSet[constructs.InterfaceInst] {
	return d.instances.Readonly()
}

func (d *interfaceDeclImp) AddInstance(inst constructs.InterfaceInst) constructs.InterfaceInst {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *interfaceDeclImp) RemoveTempDeclRefs(required bool) bool {
	if !utils.IsNil(d.nest) {
		nest, changed := constructs.ResolvedTempDeclRef(d.nest, required)
		d.nest = nest.(constructs.NestType)
		return changed
	}
	return false
}

func (d *interfaceDeclImp) FindInstance(implicitTypes, instanceTypes []constructs.TypeDesc) (constructs.InterfaceInst, bool) {
	return d.instances.Enumerate().Where(func(i constructs.InterfaceInst) bool {
		return comp.Or(
			constructs.SliceComparerPend(implicitTypes, i.ImplicitTypes()),
			constructs.SliceComparerPend(instanceTypes, i.InstanceTypes()),
		) == 0
	}).First()
}

func (d *interfaceDeclImp) IsNamed() bool   { return len(d.name) > 0 }
func (d *interfaceDeclImp) IsGeneric() bool { return len(d.typeParams) > 0 }
func (d *interfaceDeclImp) IsNested() bool  { return d.nest != nil }

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
			constructs.ComparerPend(aImp.nest, bImp.nest),
		)
	}
}

func (d *interfaceDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, d.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, d.Kind(), d.Index())
	}
	if ctx.SkipDuplicates() && d.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, d.Index()).
		Add(ctx.OnlyIndex(), `package`, d.pkg).
		Add(ctx, `name`, d.name).
		Add(ctx.OnlyIndex(), `interface`, d.inter).
		AddNonZero(ctx, `loc`, d.loc).
		AddNonZeroIf(ctx, d.exported, `vis`, `exported`).
		AddNonZero(ctx.OnlyIndex(), `typeParams`, d.typeParams).
		AddNonZero(ctx.OnlyIndex(), `nest`, d.nest).
		AddNonZero(ctx.OnlyIndex(), `instances`, d.instances.ToSlice())
}

func (d *interfaceDeclImp) ToStringer(s stringer.Stringer) {
	s.Write(d.pkg.Path(), `.`)
	if !utils.IsNil(d.nest) {
		s.Write(d.nest.Name(), `:`)
	}
	s.Write(d.name).
		WriteList(`[`, `, `, `;]`, d.ImplicitTypeParams()).
		WriteList(`[`, `, `, `]`, d.typeParams).
		Write(` interface{--}`)
}

func (d *interfaceDeclImp) String() string {
	return stringer.String(d)
}
