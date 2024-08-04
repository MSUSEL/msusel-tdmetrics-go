package constructs

import (
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Field is a variable inside of a struct or class.
//
// For the abstraction, the order of the fields and
// the tags of the fields don't matter.
type Field interface {
	Construct
	_field()

	Name() string
	Type() TypeDesc
	Embedded() bool
}

type FieldArgs struct {
	Name     string
	Type     TypeDesc
	Embedded bool
}

type fieldImp struct {
	name     string
	typ      TypeDesc
	embedded bool
	index    int
}

func newField(args FieldArgs) Field {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &fieldImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (f *fieldImp) _field()            {}
func (f *fieldImp) Kind() kind.Kind    { return kind.Field }
func (f *fieldImp) setIndex(index int) { f.index = index }

func (f *fieldImp) Name() string   { return f.name }
func (f *fieldImp) Type() TypeDesc { return f.typ }
func (f *fieldImp) Embedded() bool { return f.embedded }

func (f *fieldImp) compareTo(other Construct) int {
	b := other.(*fieldImp)
	return or(
		func() int { return strings.Compare(f.name, b.name) },
		func() int { return Compare(f.typ, b.typ) },
		func() int { return boolCompare(f.embedded, b.embedded) },
	)
}

func (f *fieldImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, f.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, f.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, f.index).
		Add(ctx2, `name`, f.name).
		Add(ctx2, `type`, f.typ).
		AddNonZero(ctx2, `embedded`, f.embedded)
}
