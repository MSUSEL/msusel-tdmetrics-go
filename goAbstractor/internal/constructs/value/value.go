package value

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/declarations"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

const Kind = `value`

type Args struct {
	Package  constructs.Package
	Name     string
	Location locs.Loc
	Type     typeDescs.TypeDesc
	Const    bool
}

type valueImp struct {
	pkg     constructs.Package
	name    string
	loc     locs.Loc
	typ     typeDescs.TypeDesc
	isConst bool
	index   int
}

func New(args Args) declarations.Value {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &valueImp{
		pkg:     args.Package,
		name:    args.Name,
		loc:     args.Location,
		typ:     args.Type,
		isConst: args.Const,
	}
}

func (v *valueImp) IsDeclaration()     {}
func (v *valueImp) IsValue()           {}
func (v *valueImp) Kind() string       { return Kind }
func (v *valueImp) SetIndex(index int) { v.index = index }

func (v *valueImp) Name() string                { return v.name }
func (v *valueImp) Location() locs.Loc          { return v.loc }
func (v *valueImp) Package() constructs.Package { return v.pkg }

func (v *valueImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[declarations.Value](v, other, Comparer())
}

func Comparer() comp.Comparer[declarations.Value] {
	return func(a, b declarations.Value) int {
		aImp, bImp := a.(*valueImp), b.(*valueImp)
		return comp.Or(
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (v *valueImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, v.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, v.index).
		Add(ctx2, `package`, v.pkg).
		Add(ctx2, `name`, v.name).
		Add(ctx2, `type`, v.typ).
		AddNonZero(ctx2, `loc`, v.loc).
		AddNonZero(ctx2, `const`, v.isConst)
}
