package tempDeclRef

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

type tempDeclRefImp struct {
	constructs.ConstructCore
	pkgPath       string
	name          string
	receiver      string
	implicitTypes []constructs.TypeDesc
	instanceTypes []constructs.TypeDesc
	nest          constructs.NestType
	con           constructs.Construct
	funcType      *types.Func
}

func newTempDeclRef(args constructs.TempDeclRefArgs) constructs.TempDeclRef {
	// pkgPath may be empty for $builtin
	assert.ArgNotEmpty(`name`, args.Name)
	assert.ArgHasNoNils(`implicit types`, args.ImplicitTypes)
	assert.ArgHasNoNils(`instance types`, args.InstanceTypes)
	if len(args.ImplicitTypes) > 0 {
		assert.ArgNotNil(`nest`, args.Nest)
	}

	return &tempDeclRefImp{
		pkgPath:       args.PackagePath,
		implicitTypes: args.ImplicitTypes,
		instanceTypes: args.InstanceTypes,
		nest:          args.Nest,
		name:          args.Name,
		receiver:      args.Receiver,
		funcType:      args.FuncType,
	}
}

func (r *tempDeclRefImp) IsTempDeclRef() {}

func (r *tempDeclRefImp) Kind() kind.Kind     { return kind.TempDeclRef }
func (r *tempDeclRefImp) PackagePath() string { return r.pkgPath }
func (r *tempDeclRefImp) Name() string        { return r.name }
func (r *tempDeclRefImp) Receiver() string    { return r.receiver }

func (r *tempDeclRefImp) FuncType() *types.Func                { return r.funcType }
func (r *tempDeclRefImp) ImplicitTypes() []constructs.TypeDesc { return r.implicitTypes }
func (r *tempDeclRefImp) InstanceTypes() []constructs.TypeDesc { return r.instanceTypes }
func (r *tempDeclRefImp) Nest() constructs.NestType            { return r.nest }
func (r *tempDeclRefImp) ResolvedType() constructs.Construct   { return r.con }

func (r *tempDeclRefImp) Resolved() bool {
	return !utils.IsNil(r.con)
}

func (r *tempDeclRefImp) SetResolution(con constructs.Construct) {
	if r.con == con {
		return
	}
	assert.ArgIsNil(`resolved`, r.con)
	assert.ArgNotNil(`con`, con)
	r.con = con
}

func (r *tempDeclRefImp) RemoveTempDeclRefs(required bool) bool {
	if !utils.IsNil(r.nest) {
		nest, changed := constructs.ResolvedTempDeclRef(r.nest, required)
		r.nest = nest.(constructs.NestType)
		return changed
	}
	return false
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
			comp.DefaultPend(aImp.receiver, bImp.receiver),
			constructs.SliceComparerPend(aImp.implicitTypes, bImp.implicitTypes),
			constructs.SliceComparerPend(aImp.instanceTypes, bImp.instanceTypes),
			constructs.ComparerPend(aImp.nest, bImp.nest),
		)
	}
}

func (r *tempDeclRefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, r.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, r.Kind(), r.Index())
	}
	if ctx.SkipDead() && !r.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && r.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, r.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, r.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, r.Alive()).
		AddNonZero(ctx, `packagePath`, r.pkgPath).
		Add(ctx, `name`, r.name).
		AddNonZero(ctx, `receiver`, r.receiver).
		AddNonZero(ctx.Short(), `resolved`, r.con).
		AddNonZero(ctx.Short(), `implicitTypes`, r.implicitTypes).
		AddNonZero(ctx.Short(), `instanceTypes`, r.instanceTypes).
		AddNonZero(ctx.OnlyIndex(), `nest`, r.nest)
}

func (r *tempDeclRefImp) ToStringer(s stringer.Stringer) {
	s.Write(`decl ref `)
	if len(r.pkgPath) > 0 {
		s.Write(r.pkgPath, `.`)
	}
	if !utils.IsNil(r.receiver) {
		s.Write(r.receiver, `.`)
	}
	if !utils.IsNil(r.nest) {
		s.Write(r.nest.Name(), `:`)
	}
	s.Write(r.name).
		WriteList(`[`, `, `, `;]`, r.implicitTypes).
		WriteList(`[`, `, `, `]`, r.instanceTypes)
}

func (r *tempDeclRefImp) String() string {
	return stringer.String(r)
}
