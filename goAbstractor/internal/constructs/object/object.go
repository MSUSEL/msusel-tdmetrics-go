package object

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/objectInst"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type objectImp struct {
	constructs.ConstructCore
	realType types.Type
	pkg      constructs.Package
	name     string
	exported bool
	loc      locs.Loc

	typeParams []constructs.TypeParam
	data       constructs.StructDesc
	inter      constructs.InterfaceDesc
	nest       constructs.NestType

	methods   collections.SortedSet[constructs.Method]
	instances collections.SortedSet[constructs.ObjectInst]
}

func newObject(args constructs.ObjectArgs) constructs.Object {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`type params`, args.TypeParams)
	assert.ArgNotNil(`data`, args.Data)

	return &objectImp{
		realType:   args.RealType,
		pkg:        args.Package,
		name:       args.Name,
		exported:   args.Exported,
		loc:        args.Location,
		typeParams: args.TypeParams,
		data:       args.Data,
		nest:       args.Nest,
		methods:    sortedSet.New(method.Comparer()),
		instances:  sortedSet.New(objectInst.Comparer()),
	}
}

func (d *objectImp) IsDeclaration() {}
func (d *objectImp) IsTypeDesc()    {}
func (d *objectImp) IsObject()      {}

func (d *objectImp) Kind() kind.Kind    { return kind.Object }
func (d *objectImp) GoType() types.Type { return d.realType }
func (d *objectImp) Name() string       { return d.name }
func (d *objectImp) Exported() bool     { return d.exported }
func (d *objectImp) Location() locs.Loc { return d.loc }

func (d *objectImp) Package() constructs.Package        { return d.pkg }
func (d *objectImp) Type() constructs.TypeDesc          { return d.data }
func (d *objectImp) Data() constructs.StructDesc        { return d.data }
func (d *objectImp) TypeParams() []constructs.TypeParam { return d.typeParams }
func (d *objectImp) Nest() constructs.NestType          { return d.nest }

func (d *objectImp) ImplicitTypeParams() []constructs.TypeParam {
	if d.nest == nil {
		return nil
	}
	if method, ok := d.nest.(constructs.Method); ok {
		return method.TypeParams()
	}
	panic(terror.New(`may not get ImplicitTypeParams from a non-method nesting declaration`).
		WithType(`nest`, d.nest))
}

func (d *objectImp) Instances() collections.ReadonlySortedSet[constructs.ObjectInst] {
	return d.instances.Readonly()
}

func (d *objectImp) AddInstance(inst constructs.ObjectInst) constructs.ObjectInst {
	v, _ := d.instances.TryAdd(inst)
	return v
}

func (d *objectImp) RemoveTempDeclRefs(required bool) bool {
	if !utils.IsNil(d.nest) {
		nest, changed := constructs.ResolvedTempDeclRef(d.nest, required)
		d.nest = nest.(constructs.NestType)
		return changed
	}
	return false
}

func (d *objectImp) FindInstance(implicitTypes, instanceTypes []constructs.TypeDesc) (constructs.ObjectInst, bool) {
	return d.instances.Enumerate().Where(func(i constructs.ObjectInst) bool {
		return comp.Or(
			constructs.SliceComparerPend(implicitTypes, i.ImplicitTypes()),
			constructs.SliceComparerPend(instanceTypes, i.InstanceTypes()),
		) == 0
	}).First()
}

func (d *objectImp) Methods() collections.ReadonlySortedSet[constructs.Method] {
	return d.methods.Readonly()
}

func (d *objectImp) AddMethod(met constructs.Method) constructs.Method {
	v, _ := d.methods.TryAdd(met)
	return v
}

func (d *objectImp) Interface() constructs.InterfaceDesc      { return d.inter }
func (d *objectImp) SetInterface(it constructs.InterfaceDesc) { d.inter = it }

func (d *objectImp) IsNamed() bool   { return len(d.name) > 0 }
func (d *objectImp) IsGeneric() bool { return len(d.typeParams) > 0 }
func (d *objectImp) IsNested() bool  { return d.nest != nil }

func (d *objectImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Object](d, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Object] {
	return func(a, b constructs.Object) int {
		aImp, bImp := a.(*objectImp), b.(*objectImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.typeParams, bImp.typeParams),
			constructs.ComparerPend(aImp.data, bImp.data),
			constructs.ComparerPend(aImp.nest, bImp.nest),
		)
	}
}

func (d *objectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, d.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, d.Kind(), d.Index())
	}
	if ctx.SkipDead() && !d.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && d.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, d.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, d.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, d.Alive()).
		Add(ctx.OnlyIndex(), `package`, d.pkg).
		Add(ctx, `name`, d.name).
		AddNonZero(ctx, `loc`, d.loc).
		AddNonZeroIf(ctx, d.exported, `vis`, `exported`).
		AddNonZero(ctx.OnlyIndex(), `typeParams`, d.typeParams).
		Add(ctx.OnlyIndex(), `data`, d.data).
		AddNonZero(ctx.OnlyIndex(), `instances`, d.instances.ToSlice()).
		AddNonZero(ctx.OnlyIndex(), `methods`, d.methods.ToSlice()).
		AddNonZero(ctx.OnlyIndex(), `nest`, d.nest).
		Add(ctx.OnlyIndex(), `interface`, d.inter)
}

func (d *objectImp) ToStringer(s stringer.Stringer) {
	s.Write(d.pkg.Path(), `.`)
	if !utils.IsNil(d.nest) {
		s.Write(d.nest.Name(), `:`)
	}
	s.Write(d.name).
		WriteList(`[`, `, `, `;]`, d.ImplicitTypeParams()).
		WriteList(`[`, `, `, `]`, d.typeParams).
		Write(` struct{--}`)
}

func (d *objectImp) String() string {
	return stringer.String(d)
}
