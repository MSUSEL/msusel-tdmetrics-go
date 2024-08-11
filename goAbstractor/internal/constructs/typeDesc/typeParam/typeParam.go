package typeParam

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `typeParam`

type TypeParam interface {
	typeDesc.TypeDesc
	_typeParam()

	Name() string
	Type() typeDesc.TypeDesc
}

type Args struct {
	Name string
	Type typeDesc.TypeDesc
}

type typeParamImp struct {
	name  string
	typ   typeDesc.TypeDesc
	index int
}

func newTypeParam(args Args) TypeParam {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &typeParamImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (t *typeParamImp) _typeParam()             {}
func (t *typeParamImp) Kind() string            { return Kind }
func (t *typeParamImp) SetIndex(index int)      { t.index = index }
func (t *typeParamImp) GoType() types.Type      { return t.typ.GoType() }
func (t *typeParamImp) Name() string            { return t.name }
func (t *typeParamImp) Type() typeDesc.TypeDesc { return t.typ }

func (t *typeParamImp) CompareTo(other constructs.Construct) int {
	return comp.Or(
		comp.Ordered[string]().Pend(t.Kind(), other.Kind()),
		Comparer().Pend(t, other.(TypeParam)),
	)
}

func Comparer() comp.Comparer[TypeParam] {
	return func(a, b TypeParam) int {
		aImp, bImp := a.(*typeParamImp), b.(*typeParamImp)
		return comp.Or(
			comp.Ordered[string]().Pend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (t *typeParamImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, t.index).
		Add(ctx2, `name`, t.name).
		Add(ctx2, `type`, t.typ)
}
