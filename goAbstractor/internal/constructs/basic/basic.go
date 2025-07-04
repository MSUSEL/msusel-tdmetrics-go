package basic

import (
	"cmp"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type basicImp struct {
	constructs.ConstructCore
	realType *types.Basic
}

func newBasic(args constructs.BasicArgs) constructs.Basic {
	assert.ArgNotNil(`real type`, args.RealType)
	rt := types.Default(types.Unalias(args.RealType)).(*types.Basic)
	switch rt.Kind() {
	case types.Invalid:
		panic(terror.New(`may not use an invalid type in basic construct`))
	case types.Byte:
		rt = types.Typ[types.Uint8]
	case types.Rune, types.UntypedRune:
		rt = types.Typ[types.Int32]
	case types.UnsafePointer:
		rt = types.Typ[types.Uintptr]
	case types.UntypedNil:
		panic(terror.New(`unexpected untyped nil in basic construct`))
	case types.Complex64, types.Complex128:
		panic(terror.New(`unexpected complex type in basic construct`))
	}
	return &basicImp{realType: rt}
}

func (t *basicImp) IsTypeDesc() {}
func (t *basicImp) IsBasic()    {}

func (t *basicImp) Kind() kind.Kind    { return kind.Basic }
func (t *basicImp) GoType() types.Type { return t.realType }

func (t *basicImp) basicKind() int {
	return int(t.realType.Kind())
}

func (t *basicImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Basic](t, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Basic] {
	return func(a, b constructs.Basic) int {
		aImp, bImp := a.(*basicImp), b.(*basicImp)
		return cmp.Compare(aImp.basicKind(), bImp.basicKind())
	}
}

func (t *basicImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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
	if ctx.IsDebugKindIncluded() || ctx.IsDebugIndexIncluded() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, t.Kind()).
			AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, t.Index()).
			AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, t.Alive()).
			Add(ctx, `name`, t.realType.Name())
	}
	return jsonify.New(ctx, t.realType.Name())
}

func (t *basicImp) ToStringer(s stringer.Stringer) { s.Write(t.realType.Name()) }

func (t *basicImp) String() string { return t.realType.Name() }
