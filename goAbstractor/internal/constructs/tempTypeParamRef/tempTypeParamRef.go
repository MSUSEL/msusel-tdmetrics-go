package tempTypeParamRef

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type tempTypeParamRefImp struct {
	constructs.ConstructCore
	realType types.Type
	context  string
	name     string
	resolved constructs.TypeDesc
}

func newTempTypeParamRef(args constructs.TempTypeParamRefArgs) constructs.TempTypeParamRef {
	assert.ArgNotNil(`real type`, args.RealType)
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgNotEmpty(`context`, args.Context)

	return &tempTypeParamRefImp{
		realType: args.RealType,
		name:     args.Name,
		context:  args.Context,
	}
}

func (r *tempTypeParamRefImp) IsTypeDesc()         {}
func (r *tempTypeParamRefImp) IsTypeTypeParamRef() {}

func (r *tempTypeParamRefImp) Kind() kind.Kind    { return kind.TempTypeParamRef }
func (r *tempTypeParamRefImp) GoType() types.Type { return r.realType }
func (r *tempTypeParamRefImp) Context() string    { return r.context }
func (r *tempTypeParamRefImp) Name() string       { return r.name }

func (r *tempTypeParamRefImp) ResolvedType() constructs.TypeDesc { return r.resolved }

func (r *tempTypeParamRefImp) Resolved() bool {
	return !utils.IsNil(r.resolved)
}

func (r *tempTypeParamRefImp) SetResolution(typ constructs.TypeDesc) {
	if r.resolved == typ {
		return
	}
	assert.ArgIsNil(`resolved`, r.resolved)
	assert.ArgNotNil(`type`, typ)
	r.resolved = typ
}

func (r *tempTypeParamRefImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.TempTypeParamRef](r, other, Comparer())
}

func Comparer() comp.Comparer[constructs.TempTypeParamRef] {
	return func(a, b constructs.TempTypeParamRef) int {
		aImp, bImp := a.(*tempTypeParamRefImp), b.(*tempTypeParamRefImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.context, bImp.context),
			comp.DefaultPend(aImp.name, bImp.name),
		)
	}
}

func (r *tempTypeParamRefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, r.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, r.Kind(), r.Index())
	}
	if ctx.SkipDuplicates() && r.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, r.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, r.Index()).
		AddNonZero(ctx.Short(), `context`, r.context).
		AddNonZero(ctx.Short(), `name`, r.name).
		AddNonZero(ctx.Short(), `type`, r.resolved)
}

func (r *tempTypeParamRefImp) ToStringer(s stringer.Stringer) {
	s.Write(`ref tp `, r.context, `:`, r.name)
}

func (r *tempTypeParamRefImp) String() string {
	return stringer.String(r)
}
