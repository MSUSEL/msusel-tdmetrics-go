package tempReference

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type tempReferenceImp struct {
	realType  *types.Named
	pkgPath   string
	name      string
	index     int
	alive     bool
	instTypes []constructs.TypeDesc
	typ       constructs.TypeDesc
}

func newTempReference(args constructs.TempReferenceArgs) constructs.TempReference {
	// pkgPath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		// Implement if needed.
		assert.NotImplemented()
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &tempReferenceImp{
		realType:  args.RealType,
		pkgPath:   args.PackagePath,
		instTypes: args.InstanceTypes,
		name:      args.Name,
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

func (r *tempReferenceImp) InstanceTypes() []constructs.TypeDesc { return r.instTypes }
func (r *tempReferenceImp) ResolvedType() constructs.TypeDesc    { return r.typ }

func (r *tempReferenceImp) Resolved() bool {
	return !utils.IsNil(r.typ)
}

func (r *tempReferenceImp) SetResolution(typ constructs.TypeDesc) {
	assert.ArgNotNil(`type`, typ)
	r.typ = typ
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
			constructs.SliceComparerPend(aImp.instTypes, bImp.instTypes),
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
		Add(ctx, `packagePath`, r.pkgPath).
		Add(ctx, `name`, r.name).
		AddNonZero(ctx.Short(), `type`, r.typ).
		AddNonZero(ctx.Short(), `instanceTypes`, r.instTypes)
}

func (r *tempReferenceImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`ref `)
	buf.WriteString(r.pkgPath)
	buf.WriteString(`.`)
	buf.WriteString(r.name)
	if len(r.instTypes) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(r.instTypes...).Join(`, `))
		buf.WriteString(`]`)
	}
	return buf.String()
}
