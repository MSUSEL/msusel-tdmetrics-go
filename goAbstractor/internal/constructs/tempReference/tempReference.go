package tempReference

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type tempReferenceImp struct {
	realType      types.Type
	pkgPath       string
	name          string
	index         int
	alive         bool
	implicitTypes []constructs.TypeDesc
	instanceTypes []constructs.TypeDesc
	nest          constructs.NestType
	typ           constructs.TypeDesc
}

func newTempReference(args constructs.TempReferenceArgs) constructs.TempReference {
	// args.PackagePath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`implicit types`, args.ImplicitTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)
	if len(args.ImplicitTypes) > 0 {
		assert.ArgNotNil(`nest`, args.Nest)
	}

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		// Implement if needed.
		panic(terror.New(`not implemented`).
			With("Package Path", args.PackagePath).
			With("Name", args.Name))
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &tempReferenceImp{
		realType:      args.RealType,
		pkgPath:       args.PackagePath,
		implicitTypes: args.ImplicitTypes,
		instanceTypes: args.InstanceTypes,
		nest:          args.Nest,
		name:          args.Name,
	}
}

func (r *tempReferenceImp) IsTypeDesc()      {}
func (r *tempReferenceImp) IsTypeReference() {}

func (r *tempReferenceImp) Kind() kind.Kind     { return kind.TempReference }
func (r *tempReferenceImp) Index() int          { return r.index }
func (r *tempReferenceImp) SetIndex(index int)  { r.index = index }
func (r *tempReferenceImp) Alive() bool         { return r.alive }
func (r *tempReferenceImp) SetAlive(alive bool) { r.alive = alive }
func (r *tempReferenceImp) GoType() types.Type  { return r.realType }
func (r *tempReferenceImp) PackagePath() string { return r.pkgPath }
func (r *tempReferenceImp) Name() string        { return r.name }

func (r *tempReferenceImp) ImplicitTypes() []constructs.TypeDesc { return r.implicitTypes }
func (r *tempReferenceImp) InstanceTypes() []constructs.TypeDesc { return r.instanceTypes }
func (r *tempReferenceImp) ResolvedType() constructs.TypeDesc    { return r.typ }
func (r *tempReferenceImp) Resolved() bool                       { return !utils.IsNil(r.typ) }
func (r *tempReferenceImp) Nest() constructs.NestType            { return r.nest }

func (r *tempReferenceImp) SetResolution(typ constructs.TypeDesc) {
	if r.typ == typ {
		return
	}
	assert.ArgIsNil(`resolved`, r.typ)
	assert.ArgNotNil(`type`, typ)
	r.typ = typ
}

func (r *tempReferenceImp) RemoveTempDeclRefs(required bool) {
	if !utils.IsNil(r.nest) {
		r.nest = constructs.ResolvedTempDeclRef(r.nest, required).(constructs.NestType)
	}
}

func (r *tempReferenceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.TempReference](r, other, Comparer())
}

func Comparer() comp.Comparer[constructs.TempReference] {
	return func(a, b constructs.TempReference) int {
		aImp, bImp := a.(*tempReferenceImp), b.(*tempReferenceImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.implicitTypes, bImp.implicitTypes),
			constructs.SliceComparerPend(aImp.instanceTypes, bImp.instanceTypes),
			constructs.ComparerPend(aImp.nest, bImp.nest),
		)
	}
}

func (r *tempReferenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, r.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, r.Kind(), r.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, r.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, r.index).
		AddNonZero(ctx, `packagePath`, r.pkgPath).
		Add(ctx, `name`, r.name).
		AddNonZero(ctx.Short(), `type`, r.typ).
		AddNonZero(ctx.Short(), `implicitTypes`, r.instanceTypes).
		AddNonZero(ctx.Short(), `instanceTypes`, r.instanceTypes).
		AddNonZero(ctx.OnlyIndex(), `nest`, r.nest)
}

func (r *tempReferenceImp) ToStringer(s stringer.Stringer) {
	s.Write(`ref `)
	if len(r.pkgPath) > 0 {
		s.Write(r.pkgPath, `.`)
	}
	if !utils.IsNil(r.nest) {
		s.Write(r.nest.Name(), `:`)
	}
	s.Write(r.name).
		WriteList(`[`, `, `, `;]`, r.implicitTypes).
		WriteList(`[`, `, `, `]`, r.instanceTypes)
}

func (r *tempReferenceImp) String() string {
	return stringer.String(r)
}
