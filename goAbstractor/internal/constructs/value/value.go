package value

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type valueImp struct {
	pkg     constructs.Package
	name    string
	loc     locs.Loc
	typ     constructs.TypeDesc
	isConst bool
	index   int
}

func newValue(args constructs.ValueArgs) constructs.Value {
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
func (v *valueImp) Kind() kind.Kind    { return kind.Value }
func (v *valueImp) SetIndex(index int) { v.index = index }

func (v *valueImp) Name() string                { return v.name }
func (v *valueImp) Location() locs.Loc          { return v.loc }
func (v *valueImp) Package() constructs.Package { return v.pkg }

func (v *valueImp) TypeParams() []constructs.TypeParam {
	return nil
}

func (v *valueImp) AddInstance(inst constructs.Instance) constructs.Instance {
	// TODO: Consider making a Value not a Declaration, just a Construct.
	panic(terror.New(`may not add an instance to a value declaration`).
		With(`name`, v.name))
}

func (v *valueImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Value](v, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Value] {
	return func(a, b constructs.Value) int {
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
		AddIf(ctx, ctx.IsKindShown(), `kind`, v.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, v.index).
		Add(ctx2, `package`, v.pkg).
		Add(ctx2, `name`, v.name).
		Add(ctx2, `type`, v.typ).
		AddNonZero(ctx2, `loc`, v.loc).
		AddNonZero(ctx2, `const`, v.isConst)
}
