package project

import (
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type projectImp struct {
	constructs.Factories
	locations locs.Set
}

func New(locs locs.Set) constructs.Project {
	return &projectImp{
		Factories: newFactory(locs),
		locations: locs,
	}
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
	index = updateIndices(p.Signatures(), index)
	index = updateIndices(p.StructDescs(), index)
	index = updateIndices(p.TypeParams(), index)
	updateIndices(p.Values(), index)
}

func updateIndices[T constructs.Construct](col collections.ReadonlySortedSet[T], index int) int {
	for i := range col.Count() {
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

func (p *projectImp) ResolveObjectInterfaces() {
	objects := p.Objects()
	for i := range objects.Count() {
		p.resolveObjectInter(objects.Get(i))
	}
}

func (p *projectImp) resolveObjectInter(obj constructs.Object) {
	methods := obj.Methods()
	abstracts := make([]constructs.Abstract, methods.Count())
	for i := range methods.Count() {
		method := methods.Get(i)
		abstracts[i] = p.NewAbstract(constructs.AbstractArgs{
			Name:      method.Name(),
			Signature: method.Signature(),
		})
	}
	it := p.NewInterfaceDesc(constructs.InterfaceDescArgs{
		Abstracts: abstracts,
		Package:   obj.Package().Source(),
	})
	obj.SetInterface(it)
}

func (p *projectImp) ResolveReferences() {
	refs := p.References()
	for i := range refs.Count() {
		p.resolveReference(refs.Get(i))
	}
}

func (p *projectImp) resolveReference(ref constructs.Reference) {
	if ref.Resolved() {
		return
	}

	if _, typ, ok := p.FindType(ref.PackagePath(), ref.Name(), true); ok {

		// TODO: Handle type parameters to find instance

		ref.SetType(typ)
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

	m.AddNonZero(ctx2, `abstracts`, p.Abstracts().ToSlice()).
		AddNonZero(ctx2, `arguments`, p.Arguments().ToSlice()).
		AddNonZero(ctx2, `fields`, p.Fields().ToSlice()).
		AddNonZero(ctx2, `packages`, p.Packages().ToSlice())

	m.AddNonZero(ctx2, `interfaceDecls`, p.InterfaceDecls().ToSlice()).
		AddNonZero(ctx2, `methods`, p.Methods().ToSlice()).
		AddNonZero(ctx2, `objects`, p.Objects().ToSlice()).
		AddNonZero(ctx2, `values`, p.Values().ToSlice())

	m.AddNonZero(ctx2, `basics`, p.Basics().ToSlice()).
		AddNonZero(ctx2, `instances`, p.Instances().ToSlice()).
		AddNonZero(ctx2, `interfaceDescs`, p.InterfaceDescs().ToSlice()).
		// Don't output the p.References()
		AddNonZero(ctx2, `signatures`, p.Signatures().ToSlice()).
		AddNonZero(ctx2, `structDescs`, p.StructDescs().ToSlice()).
		AddNonZero(ctx2, `typeParams`, p.TypeParams().ToSlice())

	return m
}
