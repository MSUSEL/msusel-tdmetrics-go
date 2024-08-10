package argument

import (
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

const Kind = `argument`

// Argument is a parameter or result in a method signature.
//
// The order of the arguments matters.
type Argument interface {
	constructs.Construct
	_argument()

	Name() string
	Type() typeDesc.TypeDesc
}

type Args struct {
	Name string
	Type typeDesc.TypeDesc
}

type argumentImp struct {
	name  string
	typ   typeDesc.TypeDesc
	index int
}

func newArgument(args Args) Argument {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &argumentImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (a *argumentImp) _argument()              {}
func (a *argumentImp) Kind() string            { return Kind }
func (a *argumentImp) SetIndex(index int)      { a.index = index }
func (a *argumentImp) Name() string            { return a.name }
func (a *argumentImp) Type() typeDesc.TypeDesc { return a.typ }

func (a *argumentImp) CompareTo(other constructs.Construct) int {
	return comp.Or(
		comp.Ordered[string]().Pend(a.Kind(), other.Kind()),
		Comparer().Pend(a, other.(Argument)),
	)
}

func Comparer() comp.Comparer[Argument] {
	return func(a, b Argument) int {
		aImp, bImp := a.(*argumentImp), b.(*argumentImp)
		return comp.Or(
			comp.Ordered[string]().Pend(aImp.name, bImp.name),
			func() int { return aImp.typ.CompareTo(bImp.typ) },
		)
	}
}

func (a *argumentImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, a.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, a.index).
		Add(ctx2, `name`, a.name).
		Add(ctx2, `type`, a.typ)
}
