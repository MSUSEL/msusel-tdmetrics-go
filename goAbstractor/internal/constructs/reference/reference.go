package reference

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type referenceImp struct {
	realType *types.Named
	pkgPath  string
	name     string

	typ constructs.TypeDesc
}

func newReference(args constructs.ReferenceArgs) constructs.Reference {
	assert.ArgNotNil(`real type`, args.RealType)
	// pkgPath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	return &referenceImp{
		realType: args.RealType,
		pkgPath:  args.PackagePath,
		name:     args.Name,
	}
}

func (r *referenceImp) IsTypeDesc()         {}
func (r *referenceImp) IsReference()        {}
func (r *referenceImp) Kind() kind.Kind     { return kind.Reference }
func (r *referenceImp) GoType() types.Type  { return r.realType }
func (r *referenceImp) PackagePath() string { return r.pkgPath }
func (r *referenceImp) Name() string        { return r.name }

func (r *referenceImp) Resolved() bool {
	return !utils.IsNil(r.typ)
}

func (r *referenceImp) SetIndex(index int) {
	panic(terror.New(`do not call SetIndex on Reference`))
}

func (r *referenceImp) SetType(typ constructs.TypeDecl) {
	assert.ArgNotNil(`type`, typ)
	r.typ = typ
}

func (r *referenceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Reference](r, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Reference] {
	return func(a, b constructs.Reference) int {
		aImp, bImp := a.(*referenceImp), b.(*referenceImp)
		return comp.Or(
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.name, bImp.name),
		)
	}
}

func (r *referenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, r.Kind()).
			AddNonZero(ctx2, `packagePath`, r.pkgPath).
			Add(ctx2, `name`, r.name).
			Add(ctx2, `type`, r.typ)
	}

	if utils.IsNil(r.typ) {
		return jsonify.New(ctx2, `failed to deref `+r.pkgPath+`.`+r.name)
	}

	return jsonify.New(ctx2, r.typ)
}
