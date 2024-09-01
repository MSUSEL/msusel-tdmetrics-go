package usage

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type usageImp struct {
	index int
	alive bool
}

func newUsage(args constructs.UsageArgs) constructs.Usage {
	return &usageImp{}
}

func (u *usageImp) IsUsage() {}

func (u *usageImp) Kind() kind.Kind     { return kind.Usage }
func (u *usageImp) Index() int          { return u.index }
func (u *usageImp) SetIndex(index int)  { u.index = index }
func (u *usageImp) Alive() bool         { return u.alive }
func (u *usageImp) SetAlive(alive bool) { u.alive = alive }

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
		// TODO: Implement
		//comp.DefaultPend(aImp.name, bImp.name),
		//constructs.ComparerPend(aImp.typ, bImp.typ),
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
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, u.index)
	//Add(ctx, `name`, f.name).
	//Add(ctx.Short(), `type`, f.typ)
	// TODO: Implement
}

func (f *usageImp) String() string {

	// TODO: Implement
	return `usage`
}
