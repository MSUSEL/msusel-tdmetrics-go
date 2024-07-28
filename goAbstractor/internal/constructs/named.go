package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Named interface {
		TypeDesc
		_named()

		Name() string
		Type() TypeDesc
	}

	NamedArgs struct {
		Name string
		Type TypeDesc
	}

	namedImp struct {
		name  string
		typ   TypeDesc
		index int
	}
)

func newNamed(args NamedArgs) Named {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	return &namedImp{
		name: args.Name,
		typ:  args.Type,
	}
}

func (t *namedImp) _named()            {}
func (t *namedImp) Kind() kind.Kind    { return kind.Named }
func (t *namedImp) SetIndex(index int) { t.index = index }
func (t *namedImp) GoType() types.Type { return t.typ.GoType() }
func (t *namedImp) Name() string       { return t.name }
func (t *namedImp) Type() TypeDesc     { return t.typ }

func (t *namedImp) CompareTo(other Construct) int {
	b := other.(*namedImp)
	if cmp := strings.Compare(t.name, b.name); cmp != 0 {
		return cmp
	}
	return Compare(t.typ, b.typ)
}

func (t *namedImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, t.typ)
}

func (t *namedImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `named`).
		AddIf(ctx, ctx.IsIndexShown(), `index`, t.index).
		Add(ctx2, `name`, t.name).
		Add(ctx2, `type`, t.typ)
}
