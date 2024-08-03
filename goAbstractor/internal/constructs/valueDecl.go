package constructs

import (
	"go/types"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	ValueDecl interface {
		Declaration
		_valueDecl()
	}

	ValueDeclArgs struct {
		Package  Package
		Name     string
		Location locs.Loc
		Type     TypeDesc
		Const    bool
	}

	valueDeclImp struct {
		pkg     Package
		name    string
		loc     locs.Loc
		typ     TypeDesc
		isConst bool
		index   int
	}
)

func newValueDecl(args ValueDeclArgs) ValueDecl {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &valueDeclImp{
		pkg:     args.Package,
		name:    args.Name,
		loc:     args.Location,
		typ:     args.Type,
		isConst: args.Const,
	}
}

func (v *valueDeclImp) _valueDecl()        {}
func (v *valueDeclImp) Kind() kind.Kind    { return kind.ValueDecl }
func (v *valueDeclImp) GoType() types.Type { return v.typ.GoType() }
func (v *valueDeclImp) SetIndex(index int) { v.index = index }
func (v *valueDeclImp) Name() string       { return v.name }
func (v *valueDeclImp) Location() locs.Loc { return v.loc }
func (v *valueDeclImp) Package() Package   { return v.pkg }

func (v *valueDeclImp) CompareTo(other Construct) int {
	b := other.(*valueDeclImp)
	if cmp := Compare(v.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(v.name, b.name); cmp != 0 {
		return cmp
	}
	return Compare(v.typ, b.typ)
}

func (v *valueDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
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

func (v *valueDeclImp) Visit(vi visitor.Visitor) {
	visitor.Visit(vi, v.typ)
}
