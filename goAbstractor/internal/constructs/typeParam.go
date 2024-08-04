package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type TypeParam interface {
	TypeDesc
	_typeParam()

	Name() string
	Type() TypeDesc
}

type TypeParamArgs struct {
	Name string
	Type TypeDesc
}

type typeParamImp struct {
	name  string
	typ   TypeDesc
	index int
}

func newTypeParam(args TypeParamArgs) TypeParam {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &typeParamImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (t *typeParamImp) _typeParam()        {}
func (t *typeParamImp) Kind() kind.Kind    { return kind.TypeParam }
func (t *typeParamImp) setIndex(index int) { t.index = index }
func (t *typeParamImp) GoType() types.Type { return t.typ.GoType() }
func (t *typeParamImp) Name() string       { return t.name }
func (t *typeParamImp) Type() TypeDesc     { return t.typ }

func (t *typeParamImp) compareTo(other Construct) int {
	b := other.(*typeParamImp)
	return or(
		func() int { return strings.Compare(t.name, b.name) },
		func() int { return Compare(t.typ, b.typ) },
	)
}

func (t *typeParamImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, t.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, t.index).
		Add(ctx2, `name`, t.name).
		Add(ctx2, `type`, t.typ)
}
