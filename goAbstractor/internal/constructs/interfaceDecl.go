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
	InterfaceDecl interface {
		Declaration
		_interfaceDecl()
		Interface() Interface
	}

	InterfaceDeclArgs struct {
		Package    Package
		Name       string
		Location   locs.Loc
		Type       Interface
		TypeParams []Named
	}

	interfaceDeclImp struct {
		pkg        Package
		name       string
		loc        locs.Loc
		typ        Interface
		index      int
		typeParams []Named
	}
)

func newInterfaceDecl(args InterfaceDeclArgs) InterfaceDecl {
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`type`, args.Type)
	assert.ArgNotNil(`location`, args.Location)

	return &interfaceDeclImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		typ:        args.Type,
		typeParams: args.TypeParams,
	}
}

func (id *interfaceDeclImp) _interfaceDecl()      {}
func (id *interfaceDeclImp) Kind() kind.Kind      { return kind.InterfaceDecl }
func (id *interfaceDeclImp) GoType() types.Type   { return id.typ.GoType() }
func (id *interfaceDeclImp) SetIndex(index int)   { id.index = index }
func (id *interfaceDeclImp) Name() string         { return id.name }
func (id *interfaceDeclImp) Location() locs.Loc   { return id.loc }
func (id *interfaceDeclImp) Package() Package     { return id.pkg }
func (id *interfaceDeclImp) Interface() Interface { return id.typ }

func (id *interfaceDeclImp) CompareTo(other Construct) int {
	b := other.(*interfaceDeclImp)
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

func (id *interfaceDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, id.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, id.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, id.index).
		Add(ctx2, `package`, id.pkg).
		Add(ctx2, `name`, id.name).
		Add(ctx2, `type`, id.typ).
		AddNonZero(ctx2, `loc`, id.loc).
		AddNonZero(ctx2, `typeParams`, id.typeParams)
}

func (id *interfaceDeclImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, id.typ)
}
