package usage

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type usageImp struct {
	target constructs.Construct
	origin constructs.Usage
	index  int
	alive  bool
}

func newUsage(args constructs.UsageArgs) constructs.Usage {
	assert.ArgNotNil(`target`, args.Target)
	return &usageImp{
		target: args.Target,
		origin: args.Origin,
	}
}

func (u *usageImp) IsUsage() {}

func (u *usageImp) Kind() kind.Kind     { return kind.Usage }
func (u *usageImp) Index() int          { return u.index }
func (u *usageImp) SetIndex(index int)  { u.index = index }
func (u *usageImp) Alive() bool         { return u.alive }
func (u *usageImp) SetAlive(alive bool) { u.alive = alive }
func (u *usageImp) HasOrigin() bool     { return !utils.IsNil(u.origin) }

func (u *usageImp) Target() constructs.Construct { return u.target }
func (u *usageImp) Origin() constructs.Usage     { return u.origin }

func (u *usageImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Usage](u, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Usage] {
	return func(a, b constructs.Usage) int {
		aImp, bImp := a.(*usageImp), b.(*usageImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.target, bImp.target),
			constructs.ComparerPend(aImp.origin, bImp.origin),
		)
	}
}

func (u *usageImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, u.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, u.Kind(), u.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, u.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, u.index).
		Add(ctx.Short(), `target`, u.target).
		AddNonZero(ctx.OnlyIndex(), `origin`, u.origin)
}

func (u *usageImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`usage `)
	if !utils.IsNil(u.origin) {
		buf.WriteString(u.origin.String())
		buf.WriteString(`:`)
	}
	buf.WriteString(u.target.String())
	return buf.String()
}
