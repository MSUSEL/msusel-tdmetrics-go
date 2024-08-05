package constructs

import (
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Argument is a parameter or result in a signature.
//
// The order of the arguments matters.
type Argument interface {
	Construct
	_argument()

	Name() string
	Type() TypeDesc
}

type ArgumentArgs struct {
	Name string
	Type TypeDesc
}

type argumentImp struct {
	name  string
	typ   TypeDesc
	index int
}

func newArgument(args ArgumentArgs) Argument {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &argumentImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (a *argumentImp) _argument()         {}
func (a *argumentImp) Kind() kind.Kind    { return kind.Field }
func (a *argumentImp) setIndex(index int) { a.index = index }

func (a *argumentImp) Name() string   { return a.name }
func (a *argumentImp) Type() TypeDesc { return a.typ }

func (a *argumentImp) compareTo(other Construct) int {
	b := other.(*fieldImp)
	return or(
		func() int { return strings.Compare(a.name, b.name) },
		func() int { return Compare(a.typ, b.typ) },
	)
}

func (a *argumentImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, a.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, a.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, a.index).
		Add(ctx2, `name`, a.name).
		Add(ctx2, `type`, a.typ)
}
