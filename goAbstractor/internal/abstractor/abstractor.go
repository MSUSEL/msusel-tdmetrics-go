package abstractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

// TODO:
// - Add analytics:
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)
//   - Indicate if a method is an accessor getter or setter (single expression).

func Abstract(ps []*packages.Package, verbose bool) *constructs.Project {
	ab := &abstractor{
		verbose: verbose,
		proj:    &constructs.Project{},
		baked:   map[string]typeDesc.TypeDesc{},
	}
	ab.initialize()
	ab.abstractProject(ps)
	ab.resolveImports()
	ab.resolveReceivers()
	ab.resolveInheritance()
	ab.proj.UpdateIndices()
	return ab.proj
}

type abstractor struct {
	verbose bool
	proj    *constructs.Project
	baked   map[string]typeDesc.TypeDesc
}

func (ab *abstractor) log(format string, args ...any) {
	if ab.verbose {
		fmt.Printf(format+"\n", args...)
	}
}

func (ab *abstractor) initialize() {
	ab.log(`initialize`)
	ab.bakeAny() // Prebake the "any" (i.e. object) into the interfaces.
}

func (ab *abstractor) abstractProject(ps []*packages.Package) {
	ab.log(`abstract project`)
	packages.Visit(ps, func(src *packages.Package) bool {
		if ap := ab.abstractPackage(src); ap != nil {
			ab.proj.Packages = append(ab.proj.Packages, ap)
		}
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) *constructs.Package {
	ab.log(`|  abstract package: %s`, src.PkgPath)
	pkg := &constructs.Package{
		Path:        src.PkgPath,
		ImportPaths: utils.SortedKeys(src.Imports),
	}
	for _, f := range src.Syntax {
		ab.addFile(pkg, src, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg *constructs.Package, src *packages.Package, f *ast.File) {
	ab.log(`|  |  add file to package: %s`, src.Fset.Position(f.Name.NamePos).Filename)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			ab.addGenDecl(pkg, src, d)
		case *ast.FuncDecl:
			ab.abstractFuncDecl(pkg, src, d)
		default:
			panic(fmt.Errorf(`unexpected declaration: %s`, pos(src, decl.Pos())))
		}
	}
}

func (ab *abstractor) addGenDecl(pkg *constructs.Package, src *packages.Package, decl *ast.GenDecl) {
	isConst := decl.Tok == token.CONST
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			ab.abstractTypeSpec(pkg, src, s)
		case *ast.ValueSpec:
			ab.abstractValueSpec(pkg, src, s, isConst)
		default:
			panic(fmt.Errorf(`unexpected specification: %s`, pos(src, spec.Pos())))
		}
	}
}

func (ab *abstractor) abstractTypeSpec(pkg *constructs.Package, src *packages.Package, spec *ast.TypeSpec) {
	tv, has := src.TypesInfo.Types[spec.Type]
	if !has {
		panic(fmt.Errorf(`type specification not found in types info: %s`, pos(src, spec.Type.Pos())))
	}
	def := &constructs.TypeDef{
		Name: spec.Name.Name,
		Type: ab.convertType(tv.Type),
	}
	pkg.Types = append(pkg.Types, def)
}

func (ab *abstractor) abstractValueSpec(pkg *constructs.Package, src *packages.Package, spec *ast.ValueSpec, isConst bool) {
	for _, name := range spec.Names {
		// TODO: Need to evaluate the initial value in case
		// it has connection to another var of calls a function.

		if name.Name == `_` {
			continue
		}

		tv, has := src.TypesInfo.Defs[name]
		if !has {
			panic(fmt.Errorf(`value specification not found in types info: %s`, pos(src, spec.Type.Pos())))
		}

		def := &constructs.ValueDef{
			Name:  name.Name,
			Const: isConst,
			Type:  ab.convertType(tv.Type()),
		}
		pkg.Values = append(pkg.Values, def)
	}
}

func (ab *abstractor) abstractFuncDecl(pkg *constructs.Package, src *packages.Package, decl *ast.FuncDecl) {
	obj := src.TypesInfo.Defs[decl.Name]
	m := &constructs.Method{
		Name:      decl.Name.Name,
		Signature: ab.convertSignature(obj.Type().(*types.Signature)),
		Metrics:   metrics.New(src.Fset, decl),
	}
	ab.determineReceiver(m, src, decl)
	pkg.Methods = append(pkg.Methods, m)
}

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) resolveImports() {
	ab.log(`resolve imports`)
	for _, p := range ab.proj.Packages {
		p.Imports = make([]*constructs.Package, 0, len(p.ImportPaths))
		for i, importPath := range p.ImportPaths {
			impPackage := ab.findPackageByPath(importPath)
			if impPackage == nil {
				panic(fmt.Errorf(`import package not found for %s: %s`, p.Path, importPath))
			}
			p.Imports[i] = impPackage
		}
	}
}

func (ab *abstractor) findPackageByPath(path string) *constructs.Package {
	for _, other := range ab.proj.Packages {
		if other.Path == path {
			return other
		}
	}
	return nil
}

func (ab *abstractor) registerInterface(t *typeDesc.Interface) *typeDesc.Interface {
	return registerType(t, &ab.proj.AllInterfaces)
}

func (ab *abstractor) registerSignature(t *typeDesc.Signature) *typeDesc.Signature {
	return registerType(t, &ab.proj.AllSignatures)
}

func (ab *abstractor) registerStruct(t *typeDesc.Struct) *typeDesc.Struct {
	return registerType(t, &ab.proj.AllStructs)
}

func registerType[T typeDesc.TypeDesc](t T, s *[]T) T {
	for _, t2 := range *s {
		if reflect.DeepEqual(t, t2) {
			return t2
		}
	}
	*s = append(*s, t)
	return t
}
