package constructs

import (
	"go/types"
	"slices"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
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

func newUnion(args UnionArgs) Union {
	assert.ArgNotNil(`real type`, args.RealType)

	slices.SortFunc(args.Exact, Compare)
	slices.SortFunc(args.Approx, Compare)

	return &unionImp{
		realType: args.RealType,
		exact:    args.Exact,
		approx:   args.Approx,
	}
}

type unionImp struct {
	realType *types.Union
	exact    []TypeDesc
	approx   []TypeDesc
	index    int
}

func (t *unionImp) _union()            {}
func (t *unionImp) Kind() kind.Kind    { return kind.Union }
func (t *unionImp) SetIndex(index int) { t.index = index }
func (t *unionImp) GoType() types.Type { return t.realType }

func (t *unionImp) CompareTo(other Construct) int {
	b := other.(*unionImp)
	if cmp := CompareSlice(t.exact, b.exact); cmp != 0 {
		return cmp
	}
	return CompareSlice(t.approx, b.approx)
}

func (t *unionImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, t.exact...)
	visitor.Visit(v, t.approx...)
}

func (t *unionImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, t.Kind()).
		AddNonZero(ctx2, `exact`, t.exact).
		AddNonZero(ctx2, `approx`, t.approx)
}
