package constructs

import (
	"go/types"
	"slices"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Class interface {
		Definition
		_class()

		//TypeParams() []Named
		//Methods() collections.ReadonlyList[Method]

		addMethod(met Method) Method
		resolveInterface(proj Project, pkg Package)
	}

	ClassArgs struct {
		Package    Package
		Name       string
		Data       TypeDesc
		TypeParams []Named
	}

	classImp struct {
		pkg        Package
		name       string
		data       TypeDesc
		typeParams []Named

		methods Set[Method]
		inter   Interface
		index   int
	}
)

func newClass(args ClassArgs) Class {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`package`, args.Package)
	assert.ArgNotNil(`data`, args.Data)

	if _, ok := args.Data.(Interface); ok {
		panic(terror.New(`may not use an interface as data in a class`).
			With(`package`, args.Package.Name()).
			With(`name`, args.Name).
			With(`data`, args.Data))
	}

	return &classImp{
		pkg:     args.Package,
		name:    args.Name,
		data:    args.Data,
		methods: NewSet[Method](),
	}
}

func (c *classImp) _class()             {}
func (c *classImp) Kind() kind.Kind     { return kind.Class }
func (c *classImp) SetIndex(index int)  { c.index = index }
func (c *classImp) GoType() types.Type  { return c.data.GoType() }
func (c *classImp) Name() string        { return c.name }
func (c *classImp) Package() Package    { return c.pkg }
func (c *classImp) TypeParams() []Named { return c.typeParams }

func (c *classImp) Methods() collections.ReadonlyList[Method] {
	return c.methods.Values()
}

func (c *classImp) addMethod(met Method) Method {
	return c.methods.Insert(met)
}

func (c *classImp) SetInterface(inter Interface) {
	c.inter = inter
}

func (c *classImp) CompareTo(other Construct) int {
	b := other.(*classImp)
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

func (c *classImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, c.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, c.Kind()).
		Add(ctx2, `package`, c.pkg).
		Add(ctx2, `name`, c.name).
		Add(ctx2, `data`, c.data).
		AddNonZero(ctx2, `typeParams`, c.typeParams).
		AddNonZero(ctx2, `methods`, c.methods).
		AddNonZero(ctx2, `interface`, c.inter)
}

func (c *classImp) Visit(v visitor.Visitor) bool {
	return visitor.Visit(v, c.data) &&
		visitor.Visit(v, c.typeParams...) &&
		visitor.VisitList(v, c.methods.Values()) &&
		visitor.Visit(v, c.inter)
}

func (c *classImp) resolveInterface(proj Project, pkg Package) {
	methods := c.methods.Values()
	itMethods := make([]Named, methods.Count())
	for i := range methods.Count() {
		m := methods.Get(i)
		itMethods[i] = proj.NewNamed(NamedArgs{
			Name: m.Name(),
			Type: m.Signature(),
		})
	}

	c.inter = proj.NewInterface(InterfaceArgs{
		Methods:    itMethods,
		TypeParams: slices.Clone(c.typeParams),
		Package:    pkg,
	})
}
