package field

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type fieldImp struct {
	name     string
	typ      constructs.TypeDesc
	embedded bool
	index    int
}

func newField(args constructs.FieldArgs) constructs.Field {
	// Blank name fields may be dropped from structs since we don't
	// need to pad out footprint or align fields.
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)

	return &fieldImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (f *fieldImp) IsField() {}

func (f *fieldImp) Kind() kind.Kind    { return kind.Field }
func (f *fieldImp) Index() int         { return f.index }
func (f *fieldImp) SetIndex(index int) { f.index = index }

func (f *fieldImp) Name() string              { return f.name }
func (f *fieldImp) Type() constructs.TypeDesc { return f.typ }
func (f *fieldImp) Embedded() bool            { return f.embedded }

func (f *fieldImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Field](f, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Field] {
	return func(a, b constructs.Field) int {
		aImp, bImp := a.(*fieldImp), b.(*fieldImp)
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.typ, bImp.typ),
			comp.DefaultPend(aImp.embedded, bImp.embedded),
		)
	}
}

func (f *fieldImp) RemoveTempReferences() {
	f.typ = constructs.ResolvedTempReference(f.typ)
}

func (f *fieldImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, f.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, f.Kind(), f.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, f.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, f.index).
		Add(ctx, `name`, f.name).
		Add(ctx.Short(), `type`, f.typ).
		AddNonZero(ctx, `embedded`, f.embedded)
}

func (f *fieldImp) String() string {
	return f.name + ` ` + f.typ.String()
}
