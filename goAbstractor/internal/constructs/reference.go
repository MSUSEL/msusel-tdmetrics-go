package constructs

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Reference interface {
		TypeDesc
		_reference()

		PackagePath() string
		Name() string
		Resolved() bool
		SetType(pkg Package, typ Definition)
	}

	ReferenceArgs struct {
		RealType    *types.Named
		PackagePath string
		Name        string
	}

	referenceImp struct {
		realType *types.Named
		pkgPath  string
		name     string

		pkg Package
		typ Definition
	}
)

func newReference(args ReferenceArgs) Reference {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotEmpty(`name`, args.Name)

	return &referenceImp{
		realType: args.RealType,
		pkgPath:  args.PackagePath,
		name:     args.Name,
	}
}

func (r *referenceImp) _reference()         {}
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

func (r *referenceImp) SetType(pkg Package, typ Definition) {
	assert.ArgNotNil(`type`, typ)
	r.pkg = pkg
	r.typ = typ
}

func (r *referenceImp) CompareTo(other Construct) int {
	b := other.(*referenceImp)
	if cmp := strings.Compare(r.pkgPath, b.pkgPath); cmp != 0 {
		return cmp
	}
	return strings.Compare(r.name, b.name)
}

func (r *referenceImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, r.typ)
}

func (r *referenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, r.Kind()).
			AddNonZero(ctx2, `packagePath`, r.pkgPath).
			Add(ctx2, `name`, r.name).
			AddNonZero(ctx2, `package`, r.pkg).
			Add(ctx2, `type`, r.typ)
	}

	if utils.IsNil(r.typ) {
		return jsonify.New(ctx2, `failed to deref `+r.pkgPath+`.`+r.name)
	}

	return jsonify.New(ctx2, r.typ)
}
