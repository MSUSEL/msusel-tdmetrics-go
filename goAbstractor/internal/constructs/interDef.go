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
	InterDef interface {
		Definition
		_interDef()
	}

	InterDefArgs struct {
		Package Package
		Name    string
		Type    TypeDesc
	}

	interDefImp struct {
		pkg   Package
		name  string
		typ   TypeDesc
		index int
	}
)

func newInterDef(args InterDefArgs) InterDef {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)

	return &interDefImp{
		pkg:  args.Package,
		name: args.Name,
		typ:  args.Type,
	}
}

func (id *interDefImp) _interDef()         {}
func (id *interDefImp) Kind() kind.Kind    { return kind.InterDef }
func (id *interDefImp) GoType() types.Type { return id.typ.GoType() }
func (id *interDefImp) SetIndex(index int) { id.index = index }
func (id *interDefImp) Name() string       { return id.name }
func (id *interDefImp) Package() Package   { return id.pkg }

func (id *interDefImp) CompareTo(other Construct) int {
	b := other.(*interDefImp)
	if cmp := Compare(id.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(id.name, b.name); cmp != 0 {
		return cmp
	}
	if cmp := Compare(id.typ, b.typ); cmp != 0 {
		return cmp
	}
	return 0
}

func (id *interDefImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, id.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, id.Kind()).
		Add(ctx2, `package`, id.pkg).
		Add(ctx2, `name`, id.name).
		Add(ctx2, `type`, id.typ)
}

func (id *interDefImp) Visit(v visitor.Visitor) bool {
	return visitor.Visit(v, id.typ)
}
