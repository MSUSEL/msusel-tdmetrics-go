package usage

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type usageImp struct {
	index  int
	alive  bool
	pkg    constructs.Package
	target constructs.Construct
	sel    constructs.Construct
}

func newUsage(args constructs.UsageArgs) constructs.Usage {
	return &usageImp{
		pkg:    args.Package,
		target: args.Target,
		sel:    args.Select,
	}
}

func (u *usageImp) IsUsage() {}

func (u *usageImp) Kind() kind.Kind     { return kind.Usage }
func (u *usageImp) Index() int          { return u.index }
func (u *usageImp) SetIndex(index int)  { u.index = index }
func (u *usageImp) Alive() bool         { return u.alive }
func (u *usageImp) SetAlive(alive bool) { u.alive = alive }

func (u *usageImp) Package() constructs.Package  { return u.pkg }
func (u *usageImp) Target() constructs.Construct { return u.target }
func (u *usageImp) Select() constructs.Construct { return u.sel }

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
			constructs.ComparerPend(aImp.pkg, bImp.pkg),
			constructs.ComparerPend(aImp.target, bImp.target),
			constructs.ComparerPend(aImp.sel, bImp.sel),
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
		AddNonZero(ctx.OnlyIndex(), `package`, u.pkg).
		AddNonZero(ctx.Short(), `target`, u.target).
		AddNonZero(ctx.Short(), `select`, u.sel)
}

func (f *usageImp) String() string {
	parts := []string{}
	if !utils.IsNil(f.pkg) {
		parts = append(parts, f.pkg.Path())
	}
	if !utils.IsNil(f.target) {
		parts = append(parts, f.target.String())
	}
	if !utils.IsNil(f.sel) {
		parts = append(parts, f.sel.String())
	}
	return strings.Join(parts, `.`)
}
