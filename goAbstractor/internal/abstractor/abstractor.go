package abstractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

func Abstract(ps []*packages.Package, logDepth int) constructs.Project {
	proj := constructs.NewProject()
	ab := &abstractor{
		logDepth: logDepth,
		packages: ps,
		baker:    baker.New(ps[0].Fset, proj),
		proj:     proj,
	}
	ab.initialize()

	ab.abstractProject()

	ab.log(1, `resolve imports`)
	proj.ResolveImports()

	ab.log(1, `resolve receivers`)
	proj.ResolveReceivers()

	ab.log(1, `resolve class interfaces`)
	proj.ResolveClassInterfaces()

	ab.resolveInheritance()
	ab.resolveReferences()

	// Finish and clean-up
	ab.prune()

	ab.log(1, `update indices`)
	proj.UpdateIndices()
	return proj
}

type abstractor struct {
	logDepth int
	packages []*packages.Package
	baker    baker.Baker
	proj     constructs.Project

	typeParamReplacer map[*types.TypeParam]*types.TypeParam
}

func (ab *abstractor) log(depth int, format string, args ...any) {
	if ab.logDepth >= depth {
		fmt.Printf(format+"\n", args...)
	}
}

func (ab *abstractor) initialize() {
	ab.baker.BakeBuiltin() // Prebake the builtin package.
	ab.baker.BakeAny()     // Prebake the "any" (i.e. object) into the interfaces.
}

func (ab *abstractor) abstractProject() {
	ab.log(1, `abstract project`)
	packages.Visit(ab.packages, func(src *packages.Package) bool {
		ab.abstractPackage(src)
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) constructs.Package {
	ab.log(2, `|  abstract package: %s`, src.PkgPath)
	pkg := ab.proj.NewPackage(constructs.PackageArgs{
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
	if it, ok := typ.(constructs.Interface); ok {
		ab.proj.NewInterDef(constructs.InterDefArgs{
			Package: pkg,
			Name:    spec.Name.Name,
			Type:    it,
		})
		return
	}

	// TODO: Get type params for classes
	tp := []constructs.Named{}

	ab.proj.NewClass(constructs.ClassArgs{
		Package:    pkg,
		Name:       spec.Name.Name,
		Data:       typ,
		TypeParams: tp,
	})
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
		ab.proj.NewValue(constructs.ValueArgs{
			Package: pkg,
			Name:    name.Name,
			Const:   isConst,
			Type:    typ,
		})
	}
}

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) resolveInheritance() {
	ab.log(1, `resolve inheritance`)

	obj := ab.baker.BakeAny()
	inters := ab.proj.AllInterfaces()
	for _, inter := range inters {
		obj.AddInheritors(inter)
	}
	for _, inter := range inters {
		inter.SetInheritance()
	}
}

func (ab *abstractor) resolveReferences() {
	ab.log(1, `resolve references`)
	for _, ref := range ab.proj.AllReferences() {
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
