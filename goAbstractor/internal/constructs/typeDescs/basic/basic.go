package basic

import (
	"cmp"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDescs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `basic`

// Basic is a base type (e.g. bool, int, string, float64).
//
// This does not contain complex types, those are treated as an interface.
type Basic interface {
	typeDescs.TypeDesc
	_basic()
}

type Args struct {
	RealType *types.Basic
}

type basicImp struct {
	realType *types.Basic
	index    int
}

func newBasic(args Args) Basic {
	rt := types.Unalias(types.Default(args.RealType)).(*types.Basic)
	return &basicImp{realType: rt}
}

func (t *basicImp) _basic()            {}
func (t *basicImp) Kind() string       { return Kind }
func (t *basicImp) SetIndex(index int) { t.index = index }
func (t *basicImp) GoType() types.Type { return t.realType }
func (t *basicImp) String() string     { return t.realType.Name() }

func (t *basicImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[Basic](t, other, Comparer())
}

func Comparer() comp.Comparer[Basic] {
	return func(a, b Basic) int {
		aInt := int(a.(*basicImp).realType.Info())
		bInt := int(b.(*basicImp).realType.Info())
		return cmp.Compare(aInt, bInt)
	}
}

func (t *basicImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	if ctx.IsKindShown() || ctx.IsIndexShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
			AddIf(ctx, ctx.IsIndexShown(), `index`, t.index).
			Add(ctx, `name`, t.realType.Name())
	}

	return jsonify.New(ctx, t.realType.Name())
}
