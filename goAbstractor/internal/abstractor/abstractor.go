package abstractor

import (
	"go/ast"
	"go/token"
	"go/types"
	"path/filepath"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
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
		ab.addFile(pkg, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg constructs.Package, f *ast.File) {
	ab.log.Log(`|  |  add file to package: %s`, pkg.Pos(f.Name.NamePos).Filename)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			ab.addGenDecl(pkg, d)
		case *ast.FuncDecl:
			ab.abstractFuncDecl(pkg, d)
		default:
			panic(terror.New(`unexpected declaration`).
				With(`pos`, pkg.Pos(decl.Pos())))
		}
	}
}

func (ab *abstractor) addGenDecl(pkg constructs.Package, decl *ast.GenDecl) {
	isConst := decl.Tok == token.CONST
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			ab.abstractTypeSpec(pkg, s)
		case *ast.ValueSpec:
			ab.abstractValueSpec(pkg, s, isConst)
		default:
			panic(terror.New(`unexpected specification`).
				With(`pos`, pkg.Pos(spec.Pos())))
		}
	}
}

func (ab *abstractor) abstractTypeSpec(pkg constructs.Package, spec *ast.TypeSpec) {
	tv, has := pkg.Source().TypesInfo.Types[spec.Type]
	if !has {
		panic(terror.New(`type specification not found in types info`).
			With(`pos`, pkg.Pos(spec.Pos())))
	}

	loc := ab.proj.NewLoc(spec.Pos())
	tp := ab.abstractFieldList(pkg, spec.TypeParams)
	typ := ab.convertType(tv.Type)

	if it, ok := typ.(constructs.Interface); ok {
		ab.proj.NewInterDef(constructs.InterDefArgs{
			Package:    pkg,
			Name:       spec.Name.Name,
			Type:       it,
			TypeParams: tp,
			Location:   loc,
		})
		return
	}

	ab.proj.NewClass(constructs.ClassArgs{
		Package:    pkg,
		Name:       spec.Name.Name,
		Data:       typ,
		TypeParams: tp,
		Location:   loc,
	})
}

func (ab *abstractor) abstractFieldList(pkg constructs.Package, fields *ast.FieldList) []constructs.Named {
	tp := []constructs.Named{}
	if utils.IsNil(fields) {
		for _, field := range fields.List {
			tv, has := pkg.Source().TypesInfo.Types[field.Type]
			if !has {
				panic(terror.New(`field list not found in types info`).
					With(`pos`, pkg.Pos(field.Pos())))
			}

			typ := ab.convertType(tv.Type)
			for _, name := range field.Names {
				tp = append(tp, ab.proj.NewNamed(constructs.NamedArgs{
					Name: name.Name,
					Type: typ,
				}))
			}
		}
	}
	return tp
}

func (ab *abstractor) abstractValueSpec(pkg constructs.Package, spec *ast.ValueSpec, isConst bool) {
	for _, name := range spec.Names {
		// TODO: Need to evaluate the initial value in case
		// it has connection to another var of calls a function.

		if blankName(name.Name) {
			// TODO: Could a black name assignment have a side effect?
			//       Maybe if metrics aren't nil, give it a non-blank name.
			//		 var _ = func() bool { /*bad init*/ }()
			continue
		}

		tv, has := pkg.Source().TypesInfo.Defs[name]
		if !has {
			panic(terror.New(`value specification not found in types info`).
				With(`pos`, pkg.Pos(spec.Pos())))
		}

		typ := ab.convertType(tv.Type())
		ab.proj.NewValue(constructs.ValueArgs{
			Package:  pkg,
			Name:     name.Name,
			Const:    isConst,
			Type:     typ,
			Location: ab.proj.NewLoc(spec.Pos()),
		})
	}
}

func (ab *abstractor) setTypeParamOverrides(args *types.TypeList, params *types.TypeParamList, pkg constructs.Package, decl *ast.FuncDecl) {
	count := args.Len()
	if count != params.Len() {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, pkg.Pos(decl.Pos())))
	}

	ab.typeParamReplacer = map[*types.TypeParam]*types.TypeParam{}
	for i := range count {
		ab.typeParamReplacer[args.At(i).(*types.TypeParam)] = params.At(i)
	}
}

func (ab *abstractor) clearTypeParamOverrides() {
	ab.typeParamReplacer = nil
}

func (ab *abstractor) abstractReceiver(pkg constructs.Package, decl *ast.FuncDecl) (bool, string) {
	if decl.Recv == nil || decl.Recv.NumFields() <= 0 {
		return false, ``
	}

	if decl.Recv.NumFields() != 1 {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, pkg.Pos(decl.Pos())))
	}

	noCopyRecv := false
	recv := pkg.Source().TypesInfo.Types[decl.Recv.List[0].Type].Type
	if p, ok := recv.(*types.Pointer); ok {
		noCopyRecv = true
		recv = p.Elem()
	}

	n, ok := recv.(*types.Named)
	if !ok {
		panic(terror.New(`function declaration has unexpected receiver type`).
			WithType(`receiver`, recv).
			With(`pos`, pkg.Pos(decl.Pos())))
	}
	ab.setTypeParamOverrides(n.TypeArgs(), n.TypeParams(), pkg, decl)

	recvName := n.Origin().Obj().Name()
	return noCopyRecv, recvName
}

func (ab *abstractor) abstractFuncDecl(pkg constructs.Package, decl *ast.FuncDecl) {
	info := pkg.Source().TypesInfo
	obj := info.Defs[decl.Name]

	noCopyRecv, recvName := ab.abstractReceiver(pkg, decl)
	sig := ab.convertSignature(obj.Type().(*types.Signature))
	ab.clearTypeParamOverrides()

	mets := metrics.New(pkg.Source().Fset, decl)
	loc := ab.proj.NewLoc(decl.Pos())

	name := decl.Name.Name
	if name == `init` && len(recvName) <= 0 && sig.Vacant() {
		name = `init#` + strconv.Itoa(pkg.InitCount())
	}

	ab.proj.NewMethod(constructs.MethodArgs{
		Package:    pkg,
		Name:       name,
		Signature:  sig,
		Metrics:    mets,
		NoCopyRecv: noCopyRecv,
		Receiver:   recvName,
		Location:   loc,
	})
}
