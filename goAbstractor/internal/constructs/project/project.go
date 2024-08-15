package project

import (
	"go/token"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/abstract"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/argument"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/basic"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/field"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/instance"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDecl"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/interfaceDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/method"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/object"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/packageCon"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/reference"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/structDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeParam"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/value"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type projectImp struct {
	constructs.AbstractFactory
	constructs.ArgumentFactory
	constructs.FieldFactory
	constructs.PackageFactory

	constructs.InterfaceDeclFactory
	constructs.MethodFactory
	constructs.ObjectFactory
	constructs.ValueFactory

	constructs.BasicFactory
	constructs.InstanceFactory
	constructs.InterfaceDescFactory
	constructs.ReferenceFactory
	constructs.SignatureFactory
	constructs.StructDescFactory
	constructs.TypeParamFactory

	locations locs.Set
}

func New(locs locs.Set) constructs.Project {
	return &projectImp{
		AbstractFactory: abstract.New(),
		ArgumentFactory: argument.New(),
		FieldFactory:    field.New(),
		PackageFactory:  packageCon.New(),

		InterfaceDeclFactory: interfaceDecl.New(),
		MethodFactory:        method.New(),
		ObjectFactory:        object.New(),
		ValueFactory:         value.New(),

		BasicFactory:         basic.New(),
		InstanceFactory:      instance.New(),
		InterfaceDescFactory: interfaceDesc.New(),
		ReferenceFactory:     reference.New(),
		SignatureFactory:     signature.New(),
		StructDescFactory:    structDesc.New(),
		TypeParamFactory:     typeParam.New(),

		locations: locs,
	}
}

func (p *projectImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}

func (p *projectImp) FindType(pkgPath, typeName string, panicOnNotFound bool) (constructs.Package, constructs.TypeDecl, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, nil, false
		}
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg constructs.Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	decl := pkg.FindTypeDecl(typeName)
	if decl == nil {
		if !panicOnNotFound {
			return pkg, nil, false
		}
		panic(terror.New(`failed to find type declaration for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath))
	}

	return pkg, decl, true
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	index := 1
	index = updateIndices(p.Abstracts(), index)
	index = updateIndices(p.Arguments(), index)
	index = updateIndices(p.Basics(), index)
	index = updateIndices(p.Fields(), index)
	index = updateIndices(p.Instances(), index)
	index = updateIndices(p.InterfaceDecls(), index)
	index = updateIndices(p.InterfaceDescs(), index)
	index = updateIndices(p.Methods(), index)
	index = updateIndices(p.Objects(), index)
	index = updateIndices(p.Packages(), index)
	// Don't index the p.References()
	index = updateIndices(p.StructDescs(), index)
	index = updateIndices(p.TypeParams(), index)
	updateIndices(p.Values(), index)
}

func updateIndices[T constructs.Construct](col collections.ReadonlySortedSet[T], index int) int {
	for i, count := 0, col.Count(); i < count; i++ {
		col.Get(i).SetIndex(index)
		index++
	}
	return index
}

func (p *projectImp) ResolveImports() {
	packages := p.Packages()
	for i := range packages.Count() {
		pkg := packages.Get(i)
		for _, importPath := range pkg.ImportPaths() {
			impPackage := p.FindPackageByPath(importPath)
			if impPackage == nil {
				panic(terror.New(`import package not found`).
					With(`package path`, pkg.Path).
					With(`import path`, importPath))
			}
			pkg.AddImport(impPackage)
		}
	}
}

func (p *projectImp) ResolveReceivers() {
	packages := p.Packages()
	for i := range packages.Count() {
		packages.Get(i).ResolveReceivers()
	}
}

func (p *projectImp) ResolveInheritance() {
	its := p.InterfaceDescs()
	roots := sortedSet.New[constructs.InterfaceDesc]()
	for i := range its.Count() {
		addInheritance(roots, its.Get(i))
	}
}

func addInheritance(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		if a.Implements(it) {
			// Yi <: X
			addInheritance(a.Inherits(), it)
		} else if it.Implements(a) {
			// X <: Yi
			it.AddInherits(a)
			siblings.RemoveRange(i, 1)
		} else {
			// Possible overlap, check for super-types in subtree.
			seekInherits(a.Inherits(), it)
		}
	}
}

func seekInherits(siblings collections.SortedSet[constructs.InterfaceDesc], it constructs.InterfaceDesc) {
	for i := siblings.Count() - 1; i >= 0; i-- {
		a := siblings.Get(i)
		if it.Implements(a) {
			it.AddInherits(a)
		} else {
			seekInherits(a.Inherits(), it)
		}
	}
}

func (p *projectImp) ResolveReferences() {
	refs := p.References()
	for i := range refs.Count() {
		if ref := refs.Get(i); !ref.Resolved() {
			if _, typ, ok := p.FindType(ref.PackagePath(), ref.Name(), true); ok {
				ref.SetType(typ)
			}
		}
	}
}

func (p *projectImp) FlagLocations() {
	p.locations.Reset()
	flagList(p.InterfaceDecls())
	flagList(p.Methods())
	flagList(p.Objects())
	flagList(p.Values())
}

func flagList[T constructs.Declaration](c collections.ReadonlySortedSet[T]) {
	for i := range c.Count() {
		c.Get(i).Location().Flag()
	}
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	m := jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `locs`, p.locations)

	m.AddNonZero(ctx2, `abstracts`, p.Abstracts()).
		AddNonZero(ctx2, `arguments`, p.Arguments()).
		AddNonZero(ctx2, `fields`, p.Fields()).
		AddNonZero(ctx2, `packages`, p.Packages())

	m.AddNonZero(ctx2, `interfaceDecls`, p.InterfaceDecls()).
		AddNonZero(ctx2, `methods`, p.Methods()).
		AddNonZero(ctx2, `objects`, p.Objects()).
		AddNonZero(ctx2, `values`, p.Values())

	m.AddNonZero(ctx2, `basics`, p.Basics()).
		AddNonZero(ctx2, `instances`, p.Instances()).
		AddNonZero(ctx2, `interfaceDescs`, p.InterfaceDescs()).
		// // Don't output the p.References()
		AddNonZero(ctx2, `signatures`, p.Signatures()).
		AddNonZero(ctx2, `structDescs`, p.StructDescs()).
		AddNonZero(ctx2, `typeParams`, p.TypeParams())

	return m
}
