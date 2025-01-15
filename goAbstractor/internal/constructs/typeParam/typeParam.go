package typeParam

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type typeParamImp struct {
	name  string
	typ   constructs.TypeDesc
	index int
	alive bool
}

func newTypeParam(args constructs.TypeParamArgs) constructs.TypeParam {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)

	// So that type params will match correctly when reading from types.Type
	// and from ast.Node, always use the type description for interfaces
	// and not the declarations.
	if itDecl, ok := args.Type.(constructs.InterfaceDecl); ok {
		args.Type = itDecl.Interface()
		assert.ArgNotNil(`decl.type`, args.Type)
	}

	return &typeParamImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (t *typeParamImp) IsTypeDesc()  {}
func (t *typeParamImp) IsTypeParam() {}

func (t *typeParamImp) Kind() kind.Kind     { return kind.TypeParam }
func (t *typeParamImp) Index() int          { return t.index }
func (t *typeParamImp) SetIndex(index int)  { t.index = index }
func (t *typeParamImp) Alive() bool         { return t.alive }
func (t *typeParamImp) SetAlive(alive bool) { t.alive = alive }
func (t *typeParamImp) GoType() types.Type  { return t.typ.GoType() }

func (t *typeParamImp) Name() string              { return t.name }
func (t *typeParamImp) Type() constructs.TypeDesc { return t.typ }

func (t *typeParamImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.TypeParam](t, other, Comparer())
}

func Comparer() comp.Comparer[constructs.TypeParam] {
	return func(a, b constructs.TypeParam) int {
		aImp, bImp := a.(*typeParamImp), b.(*typeParamImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
		)
	}
}

func (t *typeParamImp) RemoveTempReferences() {
	t.typ = constructs.ResolvedTempReference(t.typ)
}

func (t *typeParamImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, t.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, t.Kind(), t.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, t.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, t.index).
		Add(ctx, `name`, t.name).
		Add(ctx.Short(), `type`, t.typ)
}

func (t *typeParamImp) String() string {
	return t.name + ` ` + t.typ.String()
}
