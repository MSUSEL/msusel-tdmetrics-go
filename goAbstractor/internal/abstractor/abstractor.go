package abstractor

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

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

func Abstract(ps []*packages.Package, verbose bool) constructs.Project {
	ab := &abstractor{
		verbose: verbose,
		proj:    constructs.NewProject(),
		baked:   map[string]typeDesc.TypeDesc{},
	}
	ab.initialize()
	ab.abstractProject(ps)
	ab.resolveImports()
	ab.resolveReceivers()
	ab.resolveClasses()
	ab.resolveInheritance()
	ab.proj.UpdateIndices()
	return ab.proj
}

type abstractor struct {
	verbose bool
	proj    constructs.Project
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
			ab.proj.AppendPackage(ap)
		}
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) constructs.Package {
	ab.log(`|  abstract package: %s`, src.PkgPath)
	pkg := constructs.NewPackage(src, src.PkgPath, utils.SortedKeys(src.Imports))
	for _, f := range src.Syntax {
		ab.addFile(pkg, src, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg constructs.Package, src *packages.Package, f *ast.File) {
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

func (ab *abstractor) addGenDecl(pkg constructs.Package, src *packages.Package, decl *ast.GenDecl) {
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

func (ab *abstractor) abstractTypeSpec(pkg constructs.Package, src *packages.Package, spec *ast.TypeSpec) {
	tv, has := src.TypesInfo.Types[spec.Type]
	if !has {
		panic(fmt.Errorf(`type specification not found in types info: %s`, pos(src, spec.Type.Pos())))
	}

	typ := ab.convertType(tv.Type)
	def := constructs.NewTypeDef(spec.Name.Name, typ)
	pkg.AppendTypes(def)
}

func (ab *abstractor) abstractValueSpec(pkg constructs.Package, src *packages.Package, spec *ast.ValueSpec, isConst bool) {
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

		typ := ab.convertType(tv.Type())
		def := constructs.NewValueDef(name.Name, isConst, typ)
		pkg.AppendValues(def)
	}
}

func (ab *abstractor) abstractFuncDecl(pkg constructs.Package, src *packages.Package, decl *ast.FuncDecl) {
	obj := src.TypesInfo.Defs[decl.Name]
	sig := ab.convertSignature(obj.Type().(*types.Signature))
	m := constructs.NewMethod(decl.Name.Name, sig)
	m.SetMetrics(metrics.New(src.Fset, decl))
	ab.determineReceiver(m, src, decl)
	pkg.AppendMethods(m)
}

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) resolveImports() {
	ab.log(`resolve imports`)
	for _, p := range ab.proj.Packages() {
		imports := make([]constructs.Package, 0, len(p.ImportPaths()))
		for i, importPath := range p.ImportPaths() {
			impPackage := ab.findPackageByPath(importPath)
			if impPackage == nil {
				panic(fmt.Errorf(`import package not found for %s: %s`, p.Path(), importPath))
			}
			imports[i] = impPackage
		}
		p.SetImports(imports)
	}
}

func (ab *abstractor) findPackageByPath(path string) constructs.Package {
	for _, other := range ab.proj.Packages() {
		if other.Path() == path {
			return other
		}
	}
	return nil
}

func (ab *abstractor) resolveClasses() {
	ab.log(`resolve classes`)
	for _, pkg := range ab.proj.Packages() {
		ab.log(`|  resolve package: %s`, pkg.Source().PkgPath)
		for _, td := range pkg.Types() {
			ab.log(`|  |  resolve typeDef: %s`, td.Name)
			ab.resolveClass(pkg, td)
		}
	}
}

func (ab *abstractor) resolveClass(pkg constructs.Package, td constructs.TypeDef) {
	mTyp := []*types.Func{}
	methods := map[string]typeDesc.TypeDesc{}
	for _, m := range td.Methods() {
		s := m.Signature().GoType().(*types.Signature)
		f := types.NewFunc(token.NoPos, pkg.Source().Types, m.Name(), s)
		mTyp = append(mTyp, f)
		methods[m.Name()] = m.Signature()
	}

	tEmb := []types.Type{}
	typeParams := []typeDesc.Named{}
	// TODO: Fill parameter types for interface.

	iTyp := types.NewInterfaceType(mTyp, tEmb)
	if iTyp == nil {
		panic(fmt.Errorf(`failed to create an interface for %s.%s`, pkg.Source().PkgPath, td.Name()))
	}

	tInt := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		RealType:   iTyp,
		Methods:    methods,
		TypeParams: typeParams,
	})
	td.SetInterface(tInt)
}

func (ab *abstractor) resolveInheritance() {
	ab.log(`resolve inheritance`)
	inters := ab.proj.AllInterfaces()
	if len(inters) <= 0 {
		panic(errors.New(`expected the object interface at minimum but found no interfaces`))
	}
	obj := inters[0]
	if !obj.Equal(ab.bakeAny()) {
		panic(fmt.Errorf(`expected the first interface to be the "any" interface`))
	}
	for _, inter := range inters[1:] {
		obj.AddInheritors(inter)
	}
	for _, inter := range inters {
		inter.SetInheritance()
	}
}
