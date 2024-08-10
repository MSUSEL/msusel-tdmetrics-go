package field

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `field`

// Field is a variable inside of a struct or class.
//
// For the abstraction, the order of the fields and
// the tags of the fields don't matter.
type Field interface {
	constructs.Construct
	_field()

	Name() string
	Type() typeDesc.TypeDesc
	Embedded() bool
}

type Args struct {
	Name     string
	Type     typeDesc.TypeDesc
	Embedded bool
}

type fieldImp struct {
	name     string
	typ      typeDesc.TypeDesc
	embedded bool
	index    int
}

func newField(args Args) Field {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &fieldImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (f *fieldImp) _field()                 {}
func (f *fieldImp) Kind() string            { return Kind }
func (f *fieldImp) SetIndex(index int)      { f.index = index }
func (f *fieldImp) Name() string            { return f.name }
func (f *fieldImp) Type() typeDesc.TypeDesc { return f.typ }
func (f *fieldImp) Embedded() bool          { return f.embedded }

func (f *fieldImp) CompareTo(other constructs.Construct) int {
	return comp.Or(
		comp.Ordered[string]().Pend(f.Kind(), other.Kind()),
		Comparer().Pend(f, other.(Field)),
	)
}

func Comparer() comp.Comparer[Field] {
	return func(a, b Field) int {
		aImp, bImp := a.(*fieldImp), b.(*fieldImp)
		return comp.Or(
			comp.Ordered[string]().Pend(aImp.Kind(), bImp.Kind()),
			func() int { return aImp.typ.CompareTo(bImp.typ) },
			comp.Bool().Pend(aImp.embedded, bImp.embedded),
		)
	}
}

func (f *fieldImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, f.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, f.index).
		Add(ctx2, `name`, f.name).
		Add(ctx2, `type`, f.typ).
		AddNonZero(ctx2, `embedded`, f.embedded)
}
