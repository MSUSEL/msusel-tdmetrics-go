package _trash

/*
import (
	"go/token"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Project interface {
	jsonify.Jsonable

	//==========================

	NewArgument(args ArgumentArgs) Argument
	NewBasic(args BasicArgs) Basic
	NewField(args FieldArgs) Field
	NewInstance(args InstanceArgs) Instance
	NewInterface(args InterfaceArgs) Interface
	NewMethod(args MethodArgs) Method
	NewObject(args ObjectArgs) Object
	NewPackage(args PackageArgs) Package
	NewReference(args ReferenceArgs) Reference
	NewTypeParam(args TypeParamArgs) TypeParam
	NewValue(args ValueArgs) Value
	NewLoc(pos token.Pos) locs.Loc

	//==========================

	Objects() collections.ReadonlyList[Object]
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
	FlagLocations()
}

type projectImp struct {
	allArguments  Set[Argument]
	allBasics     Set[Basic]
	allFields     Set[Field]
	allInstances  Set[Instance]
	allInterfaces Set[Interface]
	allMethods    Set[Method]
	allObjects    Set[Object]
	allPackages   Set[Package]
	allReferences Set[Reference]
	allTypeParams Set[TypeParam]
	allValues     Set[Value]
	locations     locs.Set
}

func NewProject(locs locs.Set) Project {
	return &projectImp{
		allArguments:  NewSet[Argument](),
		allBasics:     NewSet[Basic](),
		allFields:     NewSet[Field](),
		allInstances:  NewSet[Instance](),
		allInterfaces: NewSet[Interface](),
		allMethods:    NewSet[Method](),
		allObjects:    NewSet[Object](),
		allPackages:   NewSet[Package](),
		allReferences: NewSet[Reference](),
		allTypeParams: NewSet[TypeParam](),
		allValues:     NewSet[Value](),
		locations:     locs,
	}
}

//==================================================================

func (p *projectImp) NewArgument(args ArgumentArgs) Argument {
	return p.allArguments.Insert(newArgument(args))
}

func (p *projectImp) NewBasic(args BasicArgs) Basic {
	return p.allBasics.Insert(newBasic(args))
}

func (p *projectImp) NewField(args FieldArgs) Field {
	return p.allFields.Insert(newField(args))
}

func (p *projectImp) NewInstance(args InstanceArgs) Instance {
	return p.allInstances.Insert(newInstance(args))
}

func (p *projectImp) NewInterface(args InterfaceArgs) Interface {
	return args.Package.addInterface(p.allInterfaces.Insert(newInterface(args)))
}

func (p *projectImp) NewMethod(args MethodArgs) Method {
	return args.Package.addMethod(p.allMethods.Insert(newMethod(args)))
}

func (p *projectImp) NewObject(args ObjectArgs) Object {
	return args.Package.addObject(p.allObjects.Insert(newObject(args)))
}

func (p *projectImp) NewPackage(args PackageArgs) Package {
	return p.allPackages.Insert(newPackage(args))
}

func (p *projectImp) NewReference(args ReferenceArgs) Reference {
	return p.allReferences.Insert(newReference(args))
}

func (p *projectImp) NewTypeParam(args TypeParamArgs) TypeParam {
	return p.allTypeParams.Insert(newTypeParam(args))
}

func (p *projectImp) NewValue(args ValueArgs) Value {
	return args.Package.addValue(p.allValues.Insert(newValue(args)))
}

func (p *projectImp) NewLoc(pos token.Pos) locs.Loc {
	return p.locations.NewLoc(pos)
}

//==================================================================

func (p *projectImp) Objects() collections.ReadonlyList[Object] {
	return p.allObjects.Values()
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

	def := pkg.findDeclaration(typeName)
	if def == nil {
		if !panicOnNotFound {
			return pkg, nil, false
		}
		panic(terror.New(`failed to find type for type def reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath))
	}

	return pkg, def, true
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allArguments.setIndices(index)
	index = p.allBasics.setIndices(index)
	index = p.allFields.setIndices(index)
	index = p.allInstances.setIndices(index)
	index = p.allMethods.setIndices(index)
	index = p.allObjects.setIndices(index)
	index = p.allPackages.setIndices(index)
	// Don't index the p.allReferences
	index = p.allTypeParams.setIndices(index)
	p.allValues.setIndices(index)
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `arguments`, p.allArguments).
		AddNonZero(ctx2, `basics`, p.allBasics).
		AddNonZero(ctx2, `fields`, p.allFields).
		AddNonZero(ctx2, `instances`, p.allInstances).
		AddNonZero(ctx2, `methods`, p.allMethods).
		AddNonZero(ctx2, `objects`, p.allObjects).
		AddNonZero(ctx2, `packages`, p.allPackages).
		// Don't output the p.allReferences
		AddNonZero(ctx2, `typeParams`, p.allTypeParams).
		AddNonZero(ctx2, `values`, p.allValues).
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

func (p *projectImp) ResolveInheritance() {
	decls := p.allObjects.Values()
	roots := []Object{}
	for i := range decls.Count() {
		roots = addInheritance(roots, decls.Get(i))
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

func (p *projectImp) FlagLocations() {
	p.locations.Reset()
	flagList(p.allMethods.Values())
	flagList(p.allObjects.Values())
	flagList(p.allValues.Values())
}

func flagList[T Declaration](c collections.ReadonlyList[T]) {
	for i := range c.Count() {
		c.Get(i).Location().Flag()
	}
}
*/
