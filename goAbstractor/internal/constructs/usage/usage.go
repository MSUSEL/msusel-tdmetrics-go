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
	name      string
	instTypes []constructs.TypeDesc
	origin    constructs.Construct
	resTyp    constructs.TypeDesc
	index     int
	alive     bool
}

func newUsage(args constructs.UsageArgs) constructs.Usage {
	return &usageImp{
		pkgPath:   args.PackagePath,
		name:      args.Name,
		instTypes: args.InstanceTypes,
		origin:    args.Origin,
	}
}

func (u *usageImp) IsUsage() {}

func (u *usageImp) Kind() kind.Kind     { return kind.Usage }
func (u *usageImp) Index() int          { return u.index }
func (u *usageImp) SetIndex(index int)  { u.index = index }
func (u *usageImp) Alive() bool         { return u.alive }
func (u *usageImp) SetAlive(alive bool) { u.alive = alive }

func (u *usageImp) PackagePath() string { return u.pkgPath }
func (u *usageImp) Name() string        { return u.name }
func (u *usageImp) HasOrigin() bool     { return !utils.IsNil(u.origin) }
func (u *usageImp) Resolved() bool      { return !utils.IsNil(u.resTyp) }

func (u *usageImp) Origin() constructs.Construct         { return u.origin }
func (u *usageImp) InstanceTypes() []constructs.TypeDesc { return u.instTypes }
func (u *usageImp) ResolvedType() constructs.TypeDesc    { return u.resTyp }

func (u *usageImp) SetResolution(typ constructs.TypeDesc) {
	assert.ArgNotNil(`type`, typ)
	u.resTyp = typ
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
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.origin, bImp.origin),
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
			Add(ctx.Short(), `type`, u.resTyp).
			AddNonZero(ctx.Short(), `origin`, u.origin)
	}
	return m.
		AddNonZero(ctx, `packagePath`, u.pkgPath).
		Add(ctx, `name`, u.name).
		AddNonZero(ctx.Short(), `instanceTypes`, u.instTypes).
		AddNonZero(ctx, `origin`, u.origin)
}

func (u *usageImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`usage `)
	if !utils.IsNil(u.origin) {
		buf.WriteString(u.origin.String())
		buf.WriteString(`:`)
	}
	if len(u.pkgPath) > 0 {
		buf.WriteString(u.pkgPath)
		buf.WriteString(`.`)
	}
	buf.WriteString(u.name)
	if len(u.instTypes) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(u.instTypes...).Join(`, `))
		buf.WriteString(`]`)
	}
	return buf.String()
}
