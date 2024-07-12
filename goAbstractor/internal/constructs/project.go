package constructs

import (
	"fmt"
	"go/types"
	"slices"
	"strconv"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Project interface {
		NewBasic(typ *types.Basic) Basic
		NewBasicFromName(pkg *packages.Package, typeName string) Basic
		NewClass(name string, typ TypeDesc) Class
		NewInterface(args InterfaceArgs) Interface
		NewNamed(name string, typ TypeDesc) Named
		NewPackage(args PackageArgs) Package
		NewReference(realType *types.Named, pkgPath, name string) Reference
		NewSignature(args SignatureArgs) Signature
		NewSolid(typ types.Type, target TypeDesc, tp ...TypeDesc) Solid
		NewStruct(args StructArgs) Struct
		NewUnion(args UnionArgs) Union

		//==========================

		AllInterfaces() []Interface
		AllReferences() []Reference

		//==========================

		ToJson(ctx *jsonify.Context) jsonify.Datum
		FindPackageByPath(path string) Package
		FindTypeDef(pkgName, tdName string) (Package, TypeDef)
		Packages() []Package
		AppendPackage(pkg ...Package)

		RemoveTypes(predict func(TypeDesc) bool)
		FilterPackage(predicate func(pkg Package) bool)

		// UpdateIndices should be called after all types have been registered
		// and all packages have been processed. This will update all the index
		// fields that will be used as references in the output models.
		UpdateIndices()
	}

	projectImp struct {
		allPackages   []Package
		allBasics     *typeSet[Basic]
		allClasses    *typeSet[Class]
		allInterfaces *typeSet[Interface]
		allNamed      *typeSet[Named]
		allReferences *typeSet[Reference]
		allSignatures *typeSet[Signature]
		allSolids     *typeSet[Solid]
		allStructs    *typeSet[Struct]
		allUnions     *typeSet[Union]
	}
)

func NewProject() Project {
	return &projectImp{
		allBasics:     newTypeSet[Basic](),
		allClasses:    newTypeSet[Class](),
		allInterfaces: newTypeSet[Interface](),
		allNamed:      newTypeSet[Named](),
		allReferences: newTypeSet[Reference](),
		allSignatures: newTypeSet[Signature](),
		allSolids:     newTypeSet[Solid](),
		allStructs:    newTypeSet[Struct](),
		allUnions:     newTypeSet[Union](),
	}
}

//==================================================================

func (p *projectImp) NewBasic(typ *types.Basic) Basic {
	return p.allBasics.Insert(newBasic(typ))
}

func (p *projectImp) NewBasicFromName(pkg *packages.Package, typeName string) Basic {
	return p.allBasics.Insert(newBasicFromName(pkg, typeName))
}

func (p *projectImp) NewClass(name string, typ TypeDesc) Class {
	return p.allClasses.Insert(newClass(name, typ))
}

func (p *projectImp) NewInterface(args InterfaceArgs) Interface {
	return p.allInterfaces.Insert(newInterface(args))
}

func (p *projectImp) NewNamed(name string, typ TypeDesc) Named {
	return p.allNamed.Insert(newNamed(name, typ))
}

func (p *projectImp) NewPackage(args PackageArgs) Package {
	pkg := newPackage(args)
	p.allPackages = append(p.allPackages, pkg)
	return pkg
}

func (p *projectImp) NewReference(realType *types.Named, pkgPath, name string) Reference {
	return p.allReferences.Insert(newReference(realType, pkgPath, name))
}

func (p *projectImp) NewSignature(args SignatureArgs) Signature {
	return p.allSignatures.Insert(newSignature(args))
}

func (p *projectImp) NewSolid(realType types.Type, target TypeDesc, tp ...TypeDesc) Solid {
	return p.allSolids.Insert(newSolid(realType, target, tp...))
}

func (p *projectImp) NewStruct(args StructArgs) Struct {
	return p.allStructs.Insert(newStruct(args))
}

func (p *projectImp) NewUnion(args UnionArgs) Union {
	return p.allUnions.Insert(newUnion(args))
}

//==================================================================

func (p *projectImp) AllInterfaces() []Interface {
	return p.allInterfaces.values
}

func (p *projectImp) AllReferences() []Reference {
	return p.allReferences.values
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `packages`, p.allPackages).
		AddNonZero(ctx2, `basics`, p.allBasics.values).
		AddNonZero(ctx2, `interfaces`, p.allInterfaces.values).
		AddNonZero(ctx2, `named`, p.allNamed.values).
		// Don't output p.allReferences
		AddNonZero(ctx2, `signatures`, p.allSignatures.values).
		AddNonZero(ctx2, `solids`, p.allSolids.values).
		AddNonZero(ctx2, `structs`, p.allStructs.values).
		AddNonZero(ctx2, `unions`, p.allUnions.values)
}

func (p *projectImp) FindPackageByPath(path string) Package {
	for _, other := range p.allPackages {
		if other.Path() == path {
			return other
		}
	}
	return nil
}

func (p *projectImp) FindTypeDef(pkgPath, tdName string) (Package, TypeDef) {
	if len(pkgPath) <= 0 {
		panic(fmt.Errorf(`must provide a non-empty package path for %q`, tdName))
	}

	pkg := p.FindPackageByPath(pkgPath)
	if pkg == nil {
		names := make([]string, len(p.Packages()))
		for i, pkg := range p.Packages() {
			names[i] = strconv.Quote(pkg.Path())
		}
		fmt.Println(`Package Paths: [` + strings.Join(names, `, `) + `]`)
		panic(fmt.Errorf(`failed to find package for type def reference for %q in %q`, tdName, pkgPath))
	}

	def := pkg.FindTypeDef(tdName)
	if def == nil {
		names := make([]string, len(pkg.Types()))
		for i, td := range pkg.Types() {
			names[i] = td.Name()
		}
		fmt.Println(pkgPath + `.TypeDefs: [` + strings.Join(names, `, `) + `]`)
		panic(fmt.Errorf(`failed to find type for type def reference for %q in %q`, tdName, pkgPath))
	}

	return pkg, def
}

func (p *projectImp) String() string {
	return jsonify.ToString(p)
}

func (p *projectImp) Packages() []Package {
	return p.allPackages
}

func (p *projectImp) AppendPackage(pkg ...Package) {
	p.allPackages = append(p.allPackages, pkg...)
}

func (p *projectImp) RemoveTypes(predict func(TypeDesc) bool) {
	p.allBasics.Remove(predict)
	p.allInterfaces.Remove(predict)
	p.allNamed.Remove(predict)
	p.allReferences.Remove(predict)
	p.allSignatures.Remove(predict)
	p.allSolids.Remove(predict)
	p.allStructs.Remove(predict)
	p.allUnions.Remove(predict)
}

func (p *projectImp) FilterPackage(predicate func(pkg Package) bool) {
	p.allPackages = slices.DeleteFunc(p.allPackages, predicate)
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allBasics.SetIndices(index)
	index = p.allInterfaces.SetIndices(index)
	index = p.allNamed.SetIndices(index)
	// Don't index p.allReferences
	index = p.allSignatures.SetIndices(index)
	index = p.allSolids.SetIndices(index)
	index = p.allStructs.SetIndices(index)
	index = p.allUnions.SetIndices(index)
	for i, pkg := range p.allPackages {
		index = pkg.setIndices(i+1, index)
	}
}
