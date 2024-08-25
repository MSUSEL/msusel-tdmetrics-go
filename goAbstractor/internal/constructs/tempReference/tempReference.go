package tempReference

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type tempReferenceImp struct {
	realType      *types.Named
	pkgPath       string
	name          string
	instanceTypes []constructs.TypeDesc

	typ constructs.TypeDesc
}

func newTempReference(args constructs.TempReferenceArgs) constructs.TempReference {
	// pkgPath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		// TODO: Implement if needed.
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &tempReferenceImp{
		realType:      args.RealType,
		pkgPath:       args.PackagePath,
		instanceTypes: args.InstanceTypes,
		name:          args.Name,
	}
}

func (r *tempReferenceImp) IsTypeDesc()  {}
func (r *tempReferenceImp) IsReference() {}

func (r *tempReferenceImp) Kind() kind.Kind     { return kind.TempReference }
func (r *tempReferenceImp) GoType() types.Type  { return r.realType }
func (r *tempReferenceImp) PackagePath() string { return r.pkgPath }
func (r *tempReferenceImp) Name() string        { return r.name }

func (r *tempReferenceImp) InstanceTypes() []constructs.TypeDesc { return r.instanceTypes }
func (r *tempReferenceImp) ResolvedType() constructs.TypeDesc    { return r.typ }

func (r *tempReferenceImp) Resolved() bool {
	return !utils.IsNil(r.typ)
}

func (r *tempReferenceImp) SetType(typ constructs.TypeDesc) {
	assert.ArgNotNil(`type`, typ)
	r.typ = typ
}

func (r *tempReferenceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.TempReference](r, other, Comparer())
}

func Comparer() comp.Comparer[constructs.TempReference] {
	return func(a, b constructs.TempReference) int {
		aImp, bImp := a.(*tempReferenceImp), b.(*tempReferenceImp)
		return comp.Or(
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.name, bImp.name),
		)
	}
}

func (r *tempReferenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, r.Kind()).
			AddNonZero(ctx2, `packagePath`, r.pkgPath).
			Add(ctx2, `name`, r.name).
			Add(ctx2, `type`, r.typ).
			AddNonZero(ctx2, `instanceTypes`, r.instanceTypes)
	}

	if utils.IsNil(r.typ) {
		return jsonify.New(ctx2, `failed to deref `+r.pkgPath+`.`+r.name)
	}

	return jsonify.New(ctx2, r.typ)
}

func (r *tempReferenceImp) String() string {
	return r.pkgPath + `:` + r.name
}
