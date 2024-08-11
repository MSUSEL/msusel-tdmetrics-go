package reference

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `reference`

type Reference interface {
	typeDescs.TypeDesc
	_reference()

	PackagePath() string
	Name() string
	Resolved() bool

	SetType(typ declarations.Declaration)
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

	typ declarations.Declaration
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

func (r *referenceImp) SetType(typ declarations.Declaration) {
	assert.ArgNotNil(`type`, typ)
	if _, ok := typ.(value.Value); ok {
		panic(terror.New(`may not use a value declaration as a reference target`).
			With(`declaration`, typ))
	}
	r.typ = typ
}

func (r *referenceImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[Reference](r, other, Comparer())
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
