package constructs

import (
	"fmt"
	"strconv"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
)

type (
	Project interface {
		jsonify.Jsonable

		//==========================

		NewBasic(args BasicArgs) Basic
		NewClass(args ClassArgs) Class
		NewInterDef(args InterDefArgs) InterDef
		NewInterface(args InterfaceArgs) Interface
		NewMethod(args MethodArgs) Method
		NewNamed(args NamedArgs) Named
		NewPackage(args PackageArgs) Package
		NewReference(args ReferenceArgs) Reference
		NewSignature(args SignatureArgs) Signature
		NewSolid(args SolidArgs) Solid
		NewStruct(args StructArgs) Struct
		NewUnion(args UnionArgs) Union
		NewValue(args ValueArgs) Value

		//==========================

		Interfaces() collections.ReadonlyList[Interface]
		Packages() collections.ReadonlyList[Package]
		References() collections.ReadonlyList[Reference]

		//==========================

		FindPackageByPath(path string) Package
		FindType(pkgPath, typeName string) (Package, TypeDesc)
		Remove(predict func(Construct) bool)

		// UpdateIndices should be called after all types have been registered
		// and all packages have been processed. This will update all the index
		// fields that will be used as references in the output models.
		UpdateIndices()
	}

	projectImp struct {
		allBasics     Set[Basic]
		allClasses    Set[Class]
		allInterDefs  Set[InterDef]
		allInterfaces Set[Interface]
		allMethods    Set[Method]
		allNamed      Set[Named]
		allPackages   Set[Package]
		allReferences Set[Reference]
		allSignatures Set[Signature]
		allSolids     Set[Solid]
		allStructs    Set[Struct]
		allUnions     Set[Union]
		allValues     Set[Value]
	}
)

func NewProject() Project {
	return &projectImp{
		allBasics:     NewSet[Basic](),
		allClasses:    NewSet[Class](),
		allInterDefs:  NewSet[InterDef](),
		allInterfaces: NewSet[Interface](),
		allMethods:    NewSet[Method](),
		allNamed:      NewSet[Named](),
		allPackages:   NewSet[Package](),
		allReferences: NewSet[Reference](),
		allSignatures: NewSet[Signature](),
		allSolids:     NewSet[Solid](),
		allStructs:    NewSet[Struct](),
		allUnions:     NewSet[Union](),
		allValues:     NewSet[Value](),
	}
}

//==================================================================

func (p *projectImp) NewBasic(args BasicArgs) Basic {
	return p.allBasics.Insert(newBasic(args))
}

func (p *projectImp) NewClass(args ClassArgs) Class {
	return p.allClasses.Insert(newClass(args))
}

func (p *projectImp) NewInterDef(args InterDefArgs) InterDef {
	return p.allInterDefs.Insert(newInterDef(args))
}

func (p *projectImp) NewInterface(args InterfaceArgs) Interface {
	return p.allInterfaces.Insert(newInterface(args))
}

func (p *projectImp) NewMethod(args MethodArgs) Method {
	return p.allMethods.Insert(newMethod(args))
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

func (p *projectImp) NewSolid(args SolidArgs) Solid {
	return p.allSolids.Insert(newSolid(args))
}

func (p *projectImp) NewStruct(args StructArgs) Struct {
	return p.allStructs.Insert(newStruct(args))
}

func (p *projectImp) NewUnion(args UnionArgs) Union {
	return p.allUnions.Insert(newUnion(args))
}

func (p *projectImp) NewValue(args ValueArgs) Value {
	return p.allValues.Insert(newValue(args))
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

func (p *projectImp) FindType(pkgPath, typeName string) (Package, TypeDesc) {
	assert.ArgNotEmpty(`pkgPath`, pkgPath)

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		names := enumerator.Select(p.Packages().Enumerate(),
			func(pkg Package) string { return strconv.Quote(pkg.Path()) }).
			Join(`, `)
		fmt.Println(`Package Paths: [` + names + `]`)
		panic(terror.New(`failed to find package for type reference`).
			With(`type name`, typeName).
			With(`package path`, pkgPath))
	}

	def := pkg.FindType(typeName)
	if def == nil {
		names := enumerator.Select(pkg.AllTypes(),
			func(td Definition) string { return td.Name() }).
			Join(`, `)
		fmt.Println(pkgPath + `.TypeDefs: [` + names + `]`)
		panic(fmt.Errorf(`failed to find type for type def reference for %q in %q`, typeName, pkgPath))
	}

	return pkg, def
}

func (p *projectImp) Remove(predict func(Construct) bool) {
	p.allBasics.Remove(predict)
	p.allClasses.Remove(predict)
	p.allInterDefs.Remove(predict)
	p.allInterfaces.Remove(predict)
	p.allMethods.Remove(predict)
	p.allNamed.Remove(predict)
	p.allPackages.Remove(predict)
	p.allReferences.Remove(predict)
	p.allSignatures.Remove(predict)
	p.allSolids.Remove(predict)
	p.allStructs.Remove(predict)
	p.allUnions.Remove(predict)
	p.allValues.Remove(predict)
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allBasics.SetIndices(index)
	index = p.allClasses.SetIndices(index)
	index = p.allInterDefs.SetIndices(index)
	index = p.allInterfaces.SetIndices(index)
	index = p.allMethods.SetIndices(index)
	index = p.allNamed.SetIndices(index)
	index = p.allPackages.SetIndices(index)
	// Don't index the p.allReferences
	index = p.allSignatures.SetIndices(index)
	index = p.allSolids.SetIndices(index)
	index = p.allStructs.SetIndices(index)
	index = p.allUnions.SetIndices(index)
	p.allValues.SetIndices(index)
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `basics`, p.allBasics).
		AddNonZero(ctx2, `classes`, p.allClasses).
		AddNonZero(ctx2, `interDefs`, p.allInterDefs).
		AddNonZero(ctx2, `interfaces`, p.allInterfaces).
		AddNonZero(ctx2, `methods`, p.allMethods).
		AddNonZero(ctx2, `named`, p.allNamed).
		AddNonZero(ctx2, `packages`, p.allPackages).
		// Don't output the p.allReferences
		AddNonZero(ctx2, `signatures`, p.allSignatures).
		AddNonZero(ctx2, `solids`, p.allSolids).
		AddNonZero(ctx2, `structs`, p.allStructs).
		AddNonZero(ctx2, `unions`, p.allUnions).
		AddNonZero(ctx2, `values`, p.allValues)
}
