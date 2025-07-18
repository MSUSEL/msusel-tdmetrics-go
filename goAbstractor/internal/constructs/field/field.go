package field

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type fieldImp struct {
	constructs.ConstructCore
	name     string
	exported bool
	typ      constructs.TypeDesc
	embedded bool
}

func newField(args constructs.FieldArgs) constructs.Field {
	// Blank name fields may be dropped from structs since we don't
	// need to pad out footprint or align fields.
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)

	return &fieldImp{
		name:     args.Name,
		exported: args.Exported,
		typ:      args.Type,
	}
}

func (f *fieldImp) IsField() {}

func (f *fieldImp) Kind() kind.Kind { return kind.Field }
func (f *fieldImp) Name() string    { return f.name }
func (f *fieldImp) Exported() bool  { return f.exported }
func (f *fieldImp) Embedded() bool  { return f.embedded }

func (f *fieldImp) Type() constructs.TypeDesc { return f.typ }

func (f *fieldImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Field](f, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Field] {
	return func(a, b constructs.Field) int {
		aImp, bImp := a.(*fieldImp), b.(*fieldImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
			comp.DefaultPend(aImp.embedded, bImp.embedded),
		)
	}
}

func (f *fieldImp) RemoveTempReferences(required bool) (changed bool) {
	f.typ, changed = constructs.ResolvedTempReference(f.typ, required)
	return
}

func (f *fieldImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, f.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, f.Kind(), f.Index())
	}
	if ctx.SkipDead() && !f.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && f.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, f.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, f.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, f.Alive()).
		Add(ctx, `name`, f.name).
		Add(ctx.Short(), `type`, f.typ).
		AddNonZeroIf(ctx, f.exported, `vis`, `exported`).
		AddNonZero(ctx, `embedded`, f.embedded)
}

func (f *fieldImp) ToStringer(s stringer.Stringer) {
	s.Write(f.name, ` `, f.typ)
}

func (f *fieldImp) String() string {
	return stringer.String(f)
}
