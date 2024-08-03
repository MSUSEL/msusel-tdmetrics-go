package constructs

import (
	"fmt"
	"go/token"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Project interface {
		jsonify.Jsonable

		//==========================

		NewBasic(args BasicArgs) Basic
		NewClassDecl(args ClassDeclArgs) ClassDecl
		NewInstance(args InstanceArgs) Instance
		NewInterfaceDecl(args InterfaceDeclArgs) InterfaceDecl
		NewInterface(args InterfaceArgs) Interface
		NewMethod(args MethodArgs) Method
		NewNamed(args NamedArgs) Named
		NewPackage(args PackageArgs) Package
		NewReference(args ReferenceArgs) Reference
		NewSignature(args SignatureArgs) Signature
		NewStruct(args StructArgs) Struct
		NewValueDecl(args ValueDeclArgs) ValueDecl
		NewLoc(pos token.Pos) locs.Loc

		//==========================

		Packages() collections.ReadonlyList[Package]
		References() collections.ReadonlyList[Reference]

		//==========================

		FindPackageByPath(path string) Package
		FindType(pkgPath, typeName string, panicOnNotFound bool) (Package, Declaration, bool)

		// UpdateIndices should be called after all types have been registered
		// and all packages have been processed. This will update all the index
		// fields that will be used as references in the output models.
		UpdateIndices()

		//==========================

		ResolveImports()
		ResolveReceivers()
		ResolveInheritance()
		ResolveReferences()
		PruneTypes()
		PrunePackages()
		FlagLocations()
	}

	projectImp struct {
		allBasics     Set[Basic]
		allClassDecls Set[ClassDecl]
		allInstances  Set[Instance]
		allInterDecls Set[InterfaceDecl]
		allInterfaces Set[Interface]
		allMethods    Set[Method]
		allNamed      Set[Named]
		allPackages   Set[Package]
		allReferences Set[Reference]
		allSignatures Set[Signature]
		allStructs    Set[Struct]
		allValueDecls Set[ValueDecl]
		locations     locs.Set
	}
)

func NewProject(locs locs.Set) Project {
	return &projectImp{
		allBasics:     NewSet[Basic](),
		allClassDecls: NewSet[ClassDecl](),
		allInstances:  NewSet[Instance](),
		allInterDecls: NewSet[InterfaceDecl](),
		allInterfaces: NewSet[Interface](),
		allMethods:    NewSet[Method](),
		allNamed:      NewSet[Named](),
		allPackages:   NewSet[Package](),
		allReferences: NewSet[Reference](),
		allSignatures: NewSet[Signature](),
		allStructs:    NewSet[Struct](),
		allValueDecls: NewSet[ValueDecl](),
		locations:     locs,
	}
}

//==================================================================

func (p *projectImp) NewBasic(args BasicArgs) Basic {
	return p.allBasics.Insert(newBasic(args))
}

func (p *projectImp) NewClassDecl(args ClassDeclArgs) ClassDecl {
	return args.Package.addClassDecl(p.allClassDecls.Insert(newClassDecl(args)))
}

func (p *projectImp) NewInstance(args InstanceArgs) Instance {
	return p.allInstances.Insert(newInstance(args))
}

func (p *projectImp) NewInterfaceDecl(args InterfaceDeclArgs) InterfaceDecl {
	return args.Package.addInterfaceDecl(p.allInterDecls.Insert(newInterfaceDecl(args)))
}

func (p *projectImp) NewInterface(args InterfaceArgs) Interface {
	return p.allInterfaces.Insert(newInterface(args))
}

func (p *projectImp) NewMethod(args MethodArgs) Method {
	return args.Package.addMethod(p.allMethods.Insert(newMethod(args)))
}

func (p *projectImp) NewNamed(args NamedArgs) Named {
	return p.allNamed.Insert(newNamed(args))
}

func (p *projectImp) NewPackage(args PackageArgs) Package {
	return p.allPackages.Insert(newPackage(args))
}

func (p *projectImp) NewReference(args ReferenceArgs) Reference {
	return p.allReferences.Insert(newReference(args))
}

func (p *projectImp) NewSignature(args SignatureArgs) Signature {
	return p.allSignatures.Insert(newSignature(args))
}

func (p *projectImp) NewStruct(args StructArgs) Struct {
	return p.allStructs.Insert(newStruct(args))
}

func (p *projectImp) NewValueDecl(args ValueDeclArgs) ValueDecl {
	return args.Package.addValueDecl(p.allValueDecls.Insert(newValueDecl(args)))
}

func (p *projectImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}

//==================================================================

func (p *projectImp) Interfaces() collections.ReadonlyList[Interface] {
	return p.allInterfaces.Values()
}

func (p *projectImp) Packages() collections.ReadonlyList[Package] {
	return p.allPackages.Values()
}

func (p *projectImp) References() collections.ReadonlyList[Reference] {
	return p.allReferences.Values()
}

//==================================================================

func (p *projectImp) FindPackageByPath(path string) Package {
	pkg, _ := p.allPackages.Values().Enumerate().
		Where(func(pkg Package) bool { return pkg.Path() == path }).
		First()
	return pkg
}

func (p *projectImp) FindType(pkgPath, typeName string, panicOnNotFound bool) (Package, Declaration, bool) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		if !panicOnNotFound {
			return nil, nil, false
		}
		names := enumerator.Select(p.allPackages.Values().Enumerate(),
			func(pkg Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		panic(terror.New(`failed to find package for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath).
			With(`existing paths`, `[`+names+`]`))
	}

	def := pkg.findType(typeName)
	if def == nil {
		if !panicOnNotFound {
			return pkg, nil, false
		}
		names := enumerator.Select(pkg.allTypes(),
			func(td Declaration) string { return td.Name() }).
			Join(`, `)
		panic(terror.New(`failed to find type for type def reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath).
			With(`type defs`, `[`+names+`]`))
	}

	return pkg, def, true
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allBasics.SetIndices(index)
	index = p.allClassDecls.SetIndices(index)
	index = p.allInstances.SetIndices(index)
	index = p.allInterDecls.SetIndices(index)
	index = p.allInterfaces.SetIndices(index)
	index = p.allMethods.SetIndices(index)
	index = p.allNamed.SetIndices(index)
	index = p.allPackages.SetIndices(index)
	// Don't index the p.allReferences
	index = p.allSignatures.SetIndices(index)
	index = p.allStructs.SetIndices(index)
	p.allValueDecls.SetIndices(index)
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `basics`, p.allBasics).
		AddNonZero(ctx2, `classDecls`, p.allClassDecls).
		AddNonZero(ctx2, `instances`, p.allInstances).
		AddNonZero(ctx2, `interfaceDecls`, p.allInterDecls).
		AddNonZero(ctx2, `interfaces`, p.allInterfaces).
		AddNonZero(ctx2, `methods`, p.allMethods).
		AddNonZero(ctx2, `named`, p.allNamed).
		AddNonZero(ctx2, `packages`, p.allPackages).
		// Don't output the p.allReferences
		AddNonZero(ctx2, `signatures`, p.allSignatures).
		AddNonZero(ctx2, `structs`, p.allStructs).
		AddNonZero(ctx2, `valueDecls`, p.allValueDecls).
		AddNonZero(ctx2, `locs`, p.locations)
}

//==================================================================

func (p *projectImp) ResolveImports() {
	packages := p.allPackages.Values()
	for i := range packages.Count() {
		pkg := packages.Get(i)
		for _, importPath := range pkg.ImportPaths() {
			impPackage := p.FindPackageByPath(importPath)
			if impPackage == nil {
				panic(terror.New(`import package not found`).
					With(`package path`, pkg.Path).
					With(`import path`, importPath))
			}
			pkg.addImport(impPackage)
		}
	}
}

func (p *projectImp) ResolveReceivers() {
	packages := p.allPackages.Values()
	for i := range packages.Count() {
		packages.Get(i).resolveReceivers()
	}
}

func printInterfaces(its []Interface, indent string) {
	for _, it := range its {
		fmt.Println(indent + it.GoType().String())
		printInterfaces(it.(*interfaceImp).inheritors, indent+`  `)
	}
}

func (p *projectImp) ResolveInheritance() {
	inters := p.allInterfaces.Values()
	roots := []Interface{}

	for i := range inters.Count() {
		roots = addInheritors(roots, inters.Get(i))
		fmt.Println(`=====================`)
		printInterfaces(roots, ``)
	}

	for i := range inters.Count() {
		inters.Get(i).setInheritance()
	}

	for i := range inters.Count() {
		inters.Get(i).sortInheritance()
	}

	classes := p.allClassDecls.Values()
	for i := range classes.Count() {
		c := classes.Get(i)
		for _, root := range roots {
			root.findImplements(c)
		}
	}
}

func (p *projectImp) ResolveReferences() {
	refs := p.References()
	for i := range refs.Count() {
		if ref := refs.Get(i); !ref.Resolved() {
			pkg, typ, _ := p.FindType(ref.PackagePath(), ref.Name(), true)
			ref.SetType(pkg, typ)
		}
	}
}

func (p *projectImp) removeTypes(predict func(Construct) bool) {
	p.allBasics.Remove(predict)
	p.allClassDecls.Remove(predict)
	p.allInstances.Remove(predict)
	p.allInterDecls.Remove(predict)
	p.allInterfaces.Remove(predict)
	p.allMethods.Remove(predict)
	p.allNamed.Remove(predict)
	p.allReferences.Remove(predict)
	p.allSignatures.Remove(predict)
	p.allStructs.Remove(predict)
	p.allValueDecls.Remove(predict)
}

func (p *projectImp) PruneTypes() {
	touched := map[Construct]bool{}

	v := visitor.New(func(value any) bool {
		if c, ok := value.(Construct); ok {
			if _, has := touched[c]; has {
				return false
			}
			touched[c] = true
		}
		return true
	})

	// Visit everything reachable from the packages.
	// Do not visit all the registered types since they are being pruned.
	visitor.VisitList(v, p.allPackages.Values())
	p.removeTypes(func(td Construct) bool {
		return !touched[td]
	})
}

func (p *projectImp) PrunePackages() {
	empty := map[Construct]bool{}
	packages := p.allPackages.Values()
	for i := range packages.Count() {
		if pkg := packages.Get(i); pkg.empty() {
			empty[pkg] = true
		}
	}

	handle := func(pkg Construct) bool {
		return empty[pkg]
	}

	p.allPackages.Remove(handle)
	allPackages := p.allPackages.Values()
	for i := range allPackages.Count() {
		allPackages.Get(i).removeImports(handle)
	}
}

func (p *projectImp) FlagLocations() {
	p.locations.Reset()
	flagList(p.allClassDecls.Values())
	flagList(p.allInterDecls.Values())
	flagList(p.allMethods.Values())
	flagList(p.allValueDecls.Values())
}

func flagList[T Declaration](c collections.ReadonlyList[T]) {
	for i := range c.Count() {
		c.Get(i).Location().Flag()
	}
}
