package constructs

import (
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	// ClassDecl is a type declaration that consists of some data and
	// a collection of functions. The data may be any type of data
	// except for an interface.
	ClassDecl interface {
		Declaration
		_classDecl()

		addMethod(met Method) Method
		addImplement(inter Interface)
	}

	ClassDeclArgs struct {
		Package    Package
		Name       string
		Location   locs.Loc
		Data       TypeDesc
		TypeParams []Named
	}

	classDeclImp struct {
		pkg        Package
		name       string
		loc        locs.Loc
		data       TypeDesc
		typeParams []Named

		methods    Set[Method]
		index      int
		implements Set[Interface]
	}
)

func newClassDecl(args ClassDeclArgs) ClassDecl {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`data`, args.Data)
	assert.ArgNotNil(`location`, args.Location)

	if _, ok := args.Data.(Interface); ok {
		panic(terror.New(`may not use an interface as data in a class`).
			With(`package`, args.Package.Name()).
			With(`name`, args.Name).
			With(`data`, args.Data))
	}

	return &classDeclImp{
		pkg:        args.Package,
		name:       args.Name,
		loc:        args.Location,
		data:       args.Data,
		typeParams: args.TypeParams,
		methods:    NewSet[Method](),
		implements: NewSet[Interface](),
	}
}

func (c *classDeclImp) _classDecl()        {}
func (c *classDeclImp) Kind() kind.Kind    { return kind.ClassDecl }
func (c *classDeclImp) SetIndex(index int) { c.index = index }
func (c *classDeclImp) GoType() types.Type { return c.data.GoType() }
func (c *classDeclImp) Location() locs.Loc { return c.loc }
func (c *classDeclImp) Name() string       { return c.name }
func (c *classDeclImp) Package() Package   { return c.pkg }

func (c *classDeclImp) addMethod(met Method) Method {
	return c.methods.Insert(met)
}

func (c *classDeclImp) addImplement(inter Interface) {
	c.implements.Insert(inter)
}

func (c *classDeclImp) CompareTo(other Construct) int {
	b := other.(*classDeclImp)
	if cmp := Compare(c.pkg, b.pkg); cmp != 0 {
		return cmp
	}
	if cmp := strings.Compare(c.name, b.name); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(c.typeParams, b.typeParams); cmp != 0 {
		return cmp
	}
	return Compare(c.data, b.data)
}

func (c *classDeclImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, c.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, c.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, c.index).
		Add(ctx2, `package`, c.pkg).
		Add(ctx2, `name`, c.name).
		Add(ctx2, `data`, c.data).
		AddNonZero(ctx2, `loc`, c.loc).
		AddNonZero(ctx2, `typeParams`, c.typeParams).
		AddNonZero(ctx2, `methods`, c.methods).
		AddNonZero(ctx2, `implements`, c.implements)
}

func (c *classDeclImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, c.data)
	visitor.Visit(v, c.typeParams...)
	visitor.VisitList(v, c.methods.Values())
	visitor.Visit(v, c.implements)
}
