package constructs

import (
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Value interface {
		Declaration
		_value()
	}

	ValueArgs struct {
		Package  Package
		Name     string
		Location locs.Loc
		Type     TypeDesc
		Const    bool
	}

	valueImp struct {
		pkg     Package
		name    string
		loc     locs.Loc
		typ     TypeDesc
		isConst bool
		index   int
	}
)

func newValue(args ValueArgs) Value {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &valueImp{
		pkg:     args.Package,
		name:    args.Name,
		loc:     args.Location,
		typ:     args.Type,
		isConst: args.Const,
	}
}

func (v *valueImp) _value()            {}
func (v *valueImp) Kind() kind.Kind    { return kind.Value }
func (v *valueImp) setIndex(index int) { v.index = index }

func (v *valueImp) Name() string       { return v.name }
func (v *valueImp) Location() locs.Loc { return v.loc }
func (v *valueImp) Package() Package   { return v.pkg }

func (v *valueImp) compareTo(other Construct) int {
	b := other.(*valueImp)
	return or(
		func() int { return Compare(v.pkg, b.pkg) },
		func() int { return strings.Compare(v.name, b.name) },
		func() int { return Compare(v.typ, b.typ) },
	)
}

func (v *valueImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, v.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, v.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, v.index).
		Add(ctx2, `package`, v.pkg).
		Add(ctx2, `name`, v.name).
		Add(ctx2, `type`, v.typ).
		AddNonZero(ctx2, `loc`, v.loc).
		AddNonZero(ctx2, `const`, v.isConst)
}

func (v *valueImp) Visit(vi visitor.Visitor) {
	visitor.Visit(vi, v.typ)
}
