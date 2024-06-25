package constructs

import (
	"fmt"
	"slices"
	"strconv"
	"strings"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Project interface {
	Types() Register
	ToJson(ctx *jsonify.Context) jsonify.Datum
	FindPackageByPath(path string) Package
	FindTypeDef(pkgName, tdName string) (Package, TypeDef)
	Packages() []Package
	AppendPackage(pkg ...Package)
	FilterPackage(predicate func(pkg Package) bool)

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()
}

type projectImp struct {
	allPackages []Package
	allTypes    Register
}

func NewProject() Project {
	return &projectImp{
		allTypes: NewRegister(),
	}
}

func (p *projectImp) Types() Register {
	return p.allTypes
}

func (p *projectImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	ctx2 := ctx.HideKind()
	return jsonify.NewMap().
		Add(ctx2, `language`, `go`).
		AddNonZero(ctx2, `types`, p.allTypes).
		AddNonZero(ctx2, `packages`, p.allPackages)
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

func (p *projectImp) FilterPackage(predicate func(pkg Package) bool) {
	p.allPackages = slices.DeleteFunc(p.allPackages, predicate)
}

func (p *projectImp) UpdateIndices() {
	// Type indices compound so that each has a unique offset.
	// The typeDefs in each package are also uniquely offset.
	index := 1
	index = p.allTypes.UpdateIndices(index)
	for i, pkg := range p.allPackages {
		index = pkg.SetIndices(i+1, index)
	}
}
