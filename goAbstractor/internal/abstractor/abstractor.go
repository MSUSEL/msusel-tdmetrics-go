package abstractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Config struct {
	Packages []*packages.Package
	Log      logger.Logger
	BasePath string
}

func Abstract(cfg Config) constructs.Project {
	assert.ArgNotNil(`log`, cfg.Log)

	basePath, err := filepath.Abs(cfg.BasePath)
	if err != nil {
		panic(terror.New(`unable to get the absolute base path`, err).
			With(`base path`, cfg.BasePath))
	}

	fs := cfg.Packages[0].Fset
	locs := locs.NewSet(fs, basePath)
	proj := constructs.NewProject(locs)
	bk := baker.New(fs, proj)

	ab := &abstractor{
		log:      cfg.Log,
		packages: cfg.Packages,
		baker:    bk,
		proj:     proj,
	}

	ab.initialize()
	ab.abstractProject()

	ab.log.Log(`resolve imports`)
	proj.ResolveImports()

	ab.log.Log(`resolve receivers`)
	proj.ResolveReceivers()

	ab.log.Log(`resolve class interfaces`)
	proj.ResolveClassInterfaces()

	ab.log.Log(`resolve inheritance`)
	proj.ResolveInheritance()

	ab.log.Log(`resolve references`)
	proj.ResolveReferences()

	ab.log.Log(`prune types`)
	proj.PruneTypes()

	ab.log.Log(`prune packages`)
	proj.PrunePackages()

	ab.log.Log(`flag locations`)
	proj.FlagLocations()

	ab.log.Log(`update indices`)
	proj.UpdateIndices()

	ab.log.Log(`done`)
	return proj
}

type abstractor struct {
	log      logger.Logger
	packages []*packages.Package
	baker    baker.Baker
	proj     constructs.Project

	typeParamReplacer map[*types.TypeParam]*types.TypeParam
}

func (ab *abstractor) initialize() {
	ab.baker.BakeBuiltin() // Prebake the builtin package.
	ab.baker.BakeAny()     // Prebake the "any" (i.e. object) into the interfaces.
}

func (ab *abstractor) abstractProject() {
	ab.log.Log(`abstract project`)
	packages.Visit(ab.packages, func(src *packages.Package) bool {
		ab.abstractPackage(src)
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) constructs.Package {
	ab.log.Log(`|  abstract package: %s`, src.PkgPath)
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
	ab.log.Log(`|  |  add file to package: %s`, src.Fset.File(f.Name.NamePos).Name())
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
			Package:  pkg,
			Name:     spec.Name.Name,
			Type:     it,
			Location: ab.proj.NewLoc(spec.Type.Pos()),
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
