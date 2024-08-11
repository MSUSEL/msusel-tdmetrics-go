package reference

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declaration"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

const Kind = `reference`

type Reference interface {
	typeDesc.TypeDesc
	_reference()

	PackagePath() string
	Name() string
	Resolved() bool

	SetType(typ declaration.Declaration)
}

type Args struct {
	RealType    *types.Named
	PackagePath string
	Name        string
}

type referenceImp struct {
	realType *types.Named
	pkgPath  string
	name     string

	typ declaration.Declaration
}

func newReference(args Args) Reference {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotEmpty(`name`, args.Name)
	return &referenceImp{
		realType: args.RealType,
		pkgPath:  args.PackagePath,
		name:     args.Name,
	}
}

func (r *referenceImp) _reference()         {}
func (r *referenceImp) Kind() string        { return Kind }
func (r *referenceImp) GoType() types.Type  { return r.realType }
func (r *referenceImp) PackagePath() string { return r.pkgPath }
func (r *referenceImp) Name() string        { return r.name }

func (r *referenceImp) Resolved() bool {
	return !utils.IsNil(r.typ)
}

func (r *referenceImp) SetIndex(index int) {
	panic(terror.New(`do not call SetIndex on Reference`))
}

func (r *referenceImp) SetType(typ declaration.Declaration) {
	assert.ArgNotNil(`type`, typ)
	r.typ = typ
}

func (r *referenceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[Reference](r, other)
}

func Comparer() comp.Comparer[Reference] {
	return func(a, b Reference) int {
		aImp, bImp := a.(*referenceImp), b.(*referenceImp)
		return comp.Or(
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.name, bImp.name),
		)
	}
}

func (r *referenceImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, r.typ)
}

func (r *referenceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind().Short()
	if ctx.IsReferenceShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
			AddNonZero(ctx2, `packagePath`, r.pkgPath).
			Add(ctx2, `name`, r.name).
			Add(ctx2, `type`, r.typ)
	}

	if utils.IsNil(r.typ) {
		return jsonify.New(ctx2, `failed to deref `+r.pkgPath+`.`+r.name)
	}

	return jsonify.New(ctx2, r.typ)
}
