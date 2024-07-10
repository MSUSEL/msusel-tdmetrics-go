package abstractor

import (
	"errors"
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func Abstract(ps []*packages.Package, logDepth int) constructs.Project {
	builtinName := `$builtin`
	builtinPkg := &packages.Package{
		PkgPath: builtinName,
		Name:    builtinName,
		Fset:    ps[0].Fset,
		Types:   types.NewPackage(builtinName, builtinName),
	}

	ab := &abstractor{
		logDepth: logDepth,
		builtin:  builtinPkg,
		packages: ps,
		proj:     constructs.NewProject(),
		baked:    map[string]any{},
	}
	ab.initialize()

	ab.abstractProject()
	ab.resolveImports()
	ab.resolveReceivers()
	ab.resolveClasses()
	ab.resolveInheritance()
	ab.resolveReferences()

	// Finish and clean-up
	ab.prune()
	ab.proj.UpdateIndices()
	return ab.proj
}

type abstractor struct {
	logDepth int
	builtin  *packages.Package
	packages []*packages.Package
	proj     constructs.Project
	baked    map[string]any

	typeParamReplacer map[*types.TypeParam]*types.TypeParam
}

func (ab *abstractor) log(depth int, format string, args ...any) {
	if ab.logDepth >= depth {
		fmt.Printf(format+"\n", args...)
	}
}

func (ab *abstractor) initialize() {
	ab.log(1, `initialize`)
	ab.bakeAny() // Prebake the "any" (i.e. object) into the interfaces.
}

func (ab *abstractor) abstractProject() {
	ab.log(1, `abstract project`)
	packages.Visit(ab.packages, func(src *packages.Package) bool {
		if ap := ab.abstractPackage(src); ap != nil {
			ab.proj.AppendPackage(ap)
		}
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) constructs.Package {
	ab.log(2, `|  abstract package: %s`, src.PkgPath)
	pkg := constructs.NewPackage(constructs.PackageArgs{
		RealPkg:     src,
		Path:        src.PkgPath,
		Name:        src.Name,
		ImportPaths: utils.SortedKeys(src.Imports),
	})
	for _, f := range src.Syntax {
		ab.addFile(pkg, src, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg constructs.Package, src *packages.Package, f *ast.File) {
	ab.log(3, `|  |  add file to package: %s`, src.Fset.Position(f.Name.NamePos).Filename)
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

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) resolveImports() {
	ab.log(1, `resolve imports`)
	for _, p := range ab.proj.Packages() {
		imports := make([]constructs.Package, len(p.ImportPaths()))
		for i, importPath := range p.ImportPaths() {
			impPackage := ab.proj.FindPackageByPath(importPath)
			if impPackage == nil {
				panic(fmt.Errorf(`import package not found for %s: %s`, p.Path(), importPath))
			}
			imports[i] = impPackage
		}
		p.SetImports(imports)
	}
}

func (ab *abstractor) resolveClasses() {
	ab.log(1, `resolve classes`)
	for _, pkg := range ab.proj.Packages() {
		ab.log(2, `|  resolve package: %s`, pkg.Source().PkgPath)
		for _, td := range pkg.Types() {
			ab.log(3, `|  |  resolve typeDef: %s`, td.Name())
			ab.resolveClass(pkg, td)
		}
	}
}

func (ab *abstractor) resolveClass(pkg constructs.Package, td constructs.TypeDef) {
	if tTyp, ok := td.Type().(constructs.Interface); ok {
		td.SetInterface(tTyp)
		return
	}

	methods := []constructs.Named{}
	for _, m := range td.Methods() {
		method := ab.proj.Types().NewNamed(m.Name(), m.Signature())
		methods = append(methods, method)
	}

	typeParams := slices.Clone(td.TypeParams())

	tInt := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
		Methods:    methods,
		TypeParams: typeParams,
		Package:    pkg.Source(),
	})
	td.SetInterface(tInt)
}

func (ab *abstractor) resolveInheritance() {
	ab.log(1, `resolve inheritance`)
	inters := ab.proj.Types().AllInterfaces()
	if len(inters) <= 0 {
		panic(errors.New(`expected the object interface at minimum but found no interfaces`))
	}

	obj := inters[0]
	if !obj.Equal(ab.bakeAny()) {
		panic(errors.New(`expected the first interface to be the "any" interface`))
	}
	for _, inter := range inters[1:] {
		obj.AddInheritors(inter)
	}
	for _, inter := range inters {
		inter.SetInheritance()
	}
}

func (ab *abstractor) resolveReferences() {
	ab.log(1, `resolve references`)
	for _, ref := range ab.proj.Types().AllReferences() {
		ab.resolveReference(ref)
	}
}

func (ab *abstractor) resolveReference(ref constructs.Reference) {
	ab.log(2, `|  resolve %s%s`, ref.PackagePath(), ref.Name())
	if len(ref.PackagePath()) > 0 {
		ref.SetType(ab.proj.FindTypeDef(ref.PackagePath(), ref.Name()))
		return
	}

	switch ref.Name() {
	case `error`:
		ref.SetType(ab.bakeError())
	case `comparable`:
		ref.SetType(ab.bakeComparable())
	default:
		panic(fmt.Errorf(`unknown reference: package=%q, name=%q`, ref.PackagePath(), ref.Name()))
	}
}
