package typeParam

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type typeParamImp struct {
	constructs.ConstructCore
	name string
	typ  constructs.TypeDesc

	loopPrevention bool
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

func (t *typeParamImp) Kind() kind.Kind    { return kind.TypeParam }
func (t *typeParamImp) GoType() types.Type { return t.typ.GoType() }

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
		if cmp := strings.Compare(aImp.name, bImp.name); cmp != 0 {
			return cmp
		}
		if aImp.loopPrevention {
			return 0
		}
		aImp.loopPrevention = true
		defer func() { aImp.loopPrevention = false }()
		return constructs.ComparerPend(aImp.typ, bImp.typ)()
	}
}

func (t *typeParamImp) RemoveTempReferences(required bool) (changed bool) {
	t.typ, changed = constructs.ResolvedTempReference(t.typ, required)
	return
}

func (t *typeParamImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, t.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, t.Kind(), t.Index())
	}
	if ctx.SkipDead() && !t.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && t.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, t.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, t.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, t.Alive()).
		Add(ctx, `name`, t.name).
		Add(ctx.Short(), `type`, t.typ)
}

func (t *typeParamImp) ToStringer(s stringer.Stringer) {
	if t.loopPrevention {
		s.Write(t.name)
		return
	}
	t.loopPrevention = true
	defer func() { t.loopPrevention = false }()
	s.Write(t.name, ` `, t.typ)
}

func (t *typeParamImp) String() string {
	return stringer.String(t)
}
