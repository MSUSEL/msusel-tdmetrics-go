package constructs

import (
	"errors"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

// Exact types are like `string|int|bool` where the type must match exactly.
// Approx types are like `~string|~int` where they type may be exact or
// an extension of the base type.
type Union interface {
	TypeDesc
	_union()
}

type UnionArgs struct {
	RealType *types.Union
	Exact    []TypeDesc
	Approx   []TypeDesc
}

func NewUnion(reg Types, args UnionArgs) Union {
	if utils.IsNil(args.RealType) {
		panic(errors.New(`must provide a real type for a union`))
	}
	return reg.RegisterUnion(&unionImp{
		realType: args.RealType,
		exact:    args.Exact,
		approx:   args.Approx,
	})
}

type unionImp struct {
	realType *types.Union
	exact    []TypeDesc
	approx   []TypeDesc
	index    int
}

func (t *unionImp) _union() {}

func (t *unionImp) Visit(v Visitor) {
	visitList(v, t.exact)
	visitList(v, t.approx)
}

func (t *unionImp) SetIndex(index int) {
	t.index = index
}

func (t *unionImp) GoType() types.Type {
	return t.realType
}

func (t *unionImp) Equal(other TypeDesc) bool {
	return equalTest(t, other, func(a, b *unionImp) bool {
		return equalList(a.exact, b.exact) &&
			equalList(a.approx, b.approx)
	})
}

func (t *unionImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `union`).
		AddNonZero(ctx2, `exact`, t.exact).
		AddNonZero(ctx2, `approx`, t.approx)
}

func (t *unionImp) String() string {
	return jsonify.ToString(t)
}
