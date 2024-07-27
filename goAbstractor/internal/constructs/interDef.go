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
	InterDef interface {
		Definition
		_interDef()
		Interface() Interface
	}

	InterDefArgs struct {
		Package    Package
		Name       string
		Location   locs.Loc
		Type       Interface
		TypeParams []Named
	}

	interDefImp struct {
		pkg        Package
		name       string
		loc        locs.Loc
		typ        Interface
		index      int
		typeParams []Named
	}
)

func newInterDef(args InterDefArgs) InterDef {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &interDefImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		typ:        args.Type,
		typeParams: args.TypeParams,
	}
}

func (id *interDefImp) _interDef()           {}
func (id *interDefImp) Kind() kind.Kind      { return kind.InterDef }
func (id *interDefImp) GoType() types.Type   { return id.typ.GoType() }
func (id *interDefImp) SetIndex(index int)   { id.index = index }
func (id *interDefImp) Name() string         { return id.name }
func (id *interDefImp) Location() locs.Loc   { return id.loc }
func (id *interDefImp) Package() Package     { return id.pkg }
func (id *interDefImp) Interface() Interface { return id.typ }

func (id *interDefImp) CompareTo(other Construct) int {
	b := other.(*interDefImp)
	if cmp := Compare(id.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(id.name, b.name); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(id.typeParams, b.typeParams); cmp != 0 {
		return cmp
	}
	return Compare(id.typ, b.typ)
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
		Add(ctx2, `type`, id.typ).
		AddNonZero(ctx2, `loc`, id.loc).
		AddNonZero(ctx2, `typeParams`, id.typeParams)
}

func (id *interDefImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, id.typ)
}
