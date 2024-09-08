package usage

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type usageImp struct {
	pkgPath   string
	target    string
	instTypes []constructs.TypeDesc
	selection string

	index int
	alive bool

	resTarget    constructs.TypeDesc
	resSelection constructs.TypeDesc
}

func newUsage(args constructs.UsageArgs) constructs.Usage {
	return &usageImp{
		pkgPath:   args.PackagePath,
		target:    args.Target,
		instTypes: args.InstanceTypes,
		selection: args.Selection,
	}
}

func (u *usageImp) IsUsage() {}

func (u *usageImp) Kind() kind.Kind     { return kind.Usage }
func (u *usageImp) Index() int          { return u.index }
func (u *usageImp) SetIndex(index int)  { u.index = index }
func (u *usageImp) Alive() bool         { return u.alive }
func (u *usageImp) SetAlive(alive bool) { u.alive = alive }

func (u *usageImp) PackagePath() string { return u.pkgPath }
func (u *usageImp) Target() string      { return u.target }
func (u *usageImp) Selection() string   { return u.selection }
func (u *usageImp) HasSelection() bool  { return len(u.selection) > 0 }

func (u *usageImp) InstanceTypes() []constructs.TypeDesc   { return u.instTypes }
func (u *usageImp) ResolvedTarget() constructs.TypeDesc    { return u.resTarget }
func (u *usageImp) ResolvedSelection() constructs.TypeDesc { return u.resSelection }
func (u *usageImp) Resolved() bool                         { return utils.IsNil(u.resTarget) }

func (u *usageImp) SetResolution(target, selection constructs.TypeDesc) {
	assert.ArgNotNil(`target`, target)
	if u.HasSelection() {
		assert.ArgNotNil(`selection`, selection)
	}
	u.resTarget = target
	u.resSelection = selection
}

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
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.target, bImp.target),
			comp.DefaultPend(aImp.selection, bImp.selection),
			constructs.SliceComparerPend(aImp.instTypes, bImp.instTypes),
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
	m := jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, u.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, u.index)
	if u.Resolved() {
		return m.
			Add(ctx.Short(), `target`, u.resTarget).
			AddNonZero(ctx.Short(), `selection`, u.resSelection)
	}
	return m.
		AddNonZero(ctx, `packagePath`, u.pkgPath).
		Add(ctx, `target`, u.target).
		AddNonZero(ctx.Short(), `instanceTypes`, u.instTypes).
		AddNonZero(ctx, `selection`, u.selection)
}

func (u *usageImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`usage `)
	if len(u.pkgPath) > 0 {
		buf.WriteString(u.pkgPath)
		buf.WriteString(`.`)
	}
	buf.WriteString(u.target)
	if len(u.instTypes) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(u.instTypes...).Join(`, `))
		buf.WriteString(`]`)
	}
	if len(u.selection) > 0 {
		buf.WriteString(`.`)
		buf.WriteString(u.selection)
	}
	return buf.String()
}
