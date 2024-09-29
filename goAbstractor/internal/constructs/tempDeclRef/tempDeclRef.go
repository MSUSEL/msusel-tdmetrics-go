package tempDeclRef

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

type tempDeclRefImp struct {
	pkgPath   string
	name      string
	index     int
	alive     bool
	instTypes []constructs.TypeDesc
	con       constructs.Construct
}

func newTempDeclRef(args constructs.TempDeclRefArgs) constructs.TempDeclRef {
	// pkgPath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)

	return &tempDeclRefImp{
		pkgPath:   args.PackagePath,
		instTypes: args.InstanceTypes,
		name:      args.Name,
	}
}

func (r *tempDeclRefImp) IsTypeMethodRef() {}

func (r *tempDeclRefImp) Kind() kind.Kind     { return kind.TempDeclRef }
func (r *tempDeclRefImp) Index() int          { return r.index }
func (r *tempDeclRefImp) SetIndex(index int)  { r.index = index }
func (r *tempDeclRefImp) Alive() bool         { return r.alive }
func (r *tempDeclRefImp) SetAlive(alive bool) { r.alive = alive }
func (r *tempDeclRefImp) PackagePath() string { return r.pkgPath }
func (r *tempDeclRefImp) Name() string        { return r.name }

func (r *tempDeclRefImp) InstanceTypes() []constructs.TypeDesc { return r.instTypes }
func (r *tempDeclRefImp) ResolvedType() constructs.Construct   { return r.con }

func (r *tempDeclRefImp) Resolved() bool {
	return !utils.IsNil(r.con)
}

func (r *tempDeclRefImp) SetResolution(con constructs.Construct) {
	assert.ArgNotNil(`con`, con)
	r.con = con
}

func (r *tempDeclRefImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.TempDeclRef](r, other, Comparer())
}

func Comparer() comp.Comparer[constructs.TempDeclRef] {
	return func(a, b constructs.TempDeclRef) int {
		aImp, bImp := a.(*tempDeclRefImp), b.(*tempDeclRefImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.pkgPath, bImp.pkgPath),
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.SliceComparerPend(aImp.instTypes, bImp.instTypes),
		)
	}
}

func (r *tempDeclRefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, r.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, r.Kind(), r.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, r.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, r.index).
		AddNonZero(ctx, `packagePath`, r.pkgPath).
		Add(ctx, `name`, r.name).
		AddNonZero(ctx.Short(), `resolved`, r.con).
		AddNonZero(ctx.Short(), `instanceTypes`, r.instTypes)
}

func (r *tempDeclRefImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`decl ref `)
	if len(r.pkgPath) > 0 {
		buf.WriteString(r.pkgPath)
		buf.WriteString(`.`)
	}
	buf.WriteString(r.name)
	if len(r.instTypes) > 0 {
		buf.WriteString(`[`)
		buf.WriteString(enumerator.Enumerate(r.instTypes...).Join(`, `))
		buf.WriteString(`]`)
	}
	return buf.String()
}
