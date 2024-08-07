package abstractor

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"path/filepath"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/metrics"
)

type Config struct {
	Packages []*packages.Package
	Logger   *log.Logger
}

func Abstract(cfg Config) constructs.Project {
	fs := cfg.Packages[0].Fset
	locs := locs.NewSet(fs)
	proj := constructs.NewProject(locs)
	bk := baker.New(fs, proj)

	ab := &abstractor{
		logger:   cfg.Logger,
		packages: cfg.Packages,
		locs:     locs,
		baker:    bk,
		proj:     proj,
	}
	ab.abstractProject()

	ab.logf(`resolve imports`)
	proj.ResolveImports()

	ab.logf(`resolve receivers`)
	proj.ResolveReceivers()

	ab.logf(`resolve inheritance`)
	proj.ResolveInheritance()

	ab.logf(`resolve references`)
	proj.ResolveReferences()

	// TODO: Improve prune to use metrics to create a dead code elimination prune.
	//ab.logf(`prune`)
	//proj.PruneTypes()
	//proj.PrunePackages()

	ab.logf(`flag locations`)
	proj.FlagLocations()

	ab.logf(`update indices`)
	proj.UpdateIndices()

	ab.logf(`done`)
	return proj
}

type abstractor struct {
	logger   *log.Logger
	packages []*packages.Package
	baker    baker.Baker
	locs     locs.Set
	proj     constructs.Project
	curPkg   constructs.Package

	typeParamReplacer map[*types.TypeParam]*types.TypeParam
}

func (ab *abstractor) logf(format string, args ...any) {
	if !utils.IsNil(ab.logger) {
		ab.logger.Printf(format, args...)
	}
}

func (ab *abstractor) pos(pos token.Pos) token.Position {
	return ab.curPkg.Source().Fset.Position(pos)
}

func (ab *abstractor) abstractProject() {
	ab.logf(`abstract project`)
	packages.Visit(ab.packages, func(src *packages.Package) bool {
		ab.abstractPackage(src)
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package) {
	ab.logf(`|  abstract package: %s`, src.PkgPath)
	ab.curPkg = ab.proj.NewPackage(constructs.PackageArgs{
		RealPkg:     src,
		Path:        src.PkgPath,
		Name:        src.Name,
		ImportPaths: utils.SortedKeys(src.Imports),
	})
	for _, f := range src.Syntax {
		ab.abstractFile(f)
	}
}

func (ab *abstractor) abstractFile(f *ast.File) {
	path := ab.pos(f.FileStart).Filename
	basePath := filepath.Base(path)
	pkgPath := ab.curPkg.Source().PkgPath
	if pkgPath != `command-line-arguments` {
		alias := filepath.ToSlash(filepath.Join(pkgPath, basePath))
		ab.locs.Alias(path, alias)
	} else {
		ab.locs.Alias(path, basePath)
	}

	ab.logf(`|  |  add file to package: %s`, basePath)
	for _, decl := range f.Decls {
		switch d := decl.(type) {
		case *ast.GenDecl:
			ab.abstractGenDecl(d)
		case *ast.FuncDecl:
			ab.abstractFuncDecl(d)
		default:
			panic(terror.New(`unexpected declaration`).
				With(`pos`, ab.pos(decl.Pos())))
		}
	}
}

func (ab *abstractor) abstractGenDecl(decl *ast.GenDecl) {
	isConst := decl.Tok == token.CONST
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			ab.abstractTypeSpec(s)
		case *ast.ValueSpec:
			ab.abstractValueSpec(s, isConst)
		default:
			panic(terror.New(`unexpected specification`).
				With(`pos`, ab.pos(spec.Pos())))
		}
	}
}

func (ab *abstractor) abstractTypeSpec(spec *ast.TypeSpec) {
	tv, has := ab.curPkg.Source().TypesInfo.Types[spec.Type]
	if !has {
		panic(terror.New(`type specification not found in types info`).
			With(`pos`, ab.pos(spec.Pos())))
	}

	loc := ab.proj.NewLoc(spec.Pos())
	tp := ab.abstractFieldList(spec.TypeParams)
	typ := ab.convertType(tv.Type)

	if it, ok := typ.(constructs.Interface); ok {
		ab.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    ab.curPkg,
			Name:       spec.Name.Name,
			Type:       it,
			TypeParams: tp,
			Location:   loc,
		})
		return
	}

	ab.proj.NewClassDecl(constructs.ClassDeclArgs{
		Package:    ab.curPkg,
		Name:       spec.Name.Name,
		Data:       typ,
		TypeParams: tp,
		Location:   loc,
	})
}

func (ab *abstractor) abstractFieldList(fields *ast.FieldList) []constructs.Named {
	ns := []constructs.Named{}
	if !utils.IsNil(fields) {
		for _, field := range fields.List {
			ns = append(ns, ab.abstractField(field)...)
		}
	}
	return ns
}

func (ab *abstractor) abstractField(field *ast.Field) []constructs.Named {
	ns := []constructs.Named{}
	if utils.IsNil(field) {
		return ns
	}

	tv, has := ab.curPkg.Source().TypesInfo.Types[field.Type]
	if !has {
		panic(terror.New(`field not found in types info`).
			With(`pos`, ab.pos(field.Pos())))
	}

	typ := ab.convertType(tv.Type)
	for _, name := range field.Names {
		named := ab.proj.NewNamed(constructs.NamedArgs{
			Name: name.Name,
			Type: typ,
		})
		ns = append(ns, named)
	}
	return ns
}

func (ab *abstractor) abstractValueSpec(spec *ast.ValueSpec, isConst bool) {
	for _, name := range spec.Names {
		// TODO: Need to evaluate the initial value in case
		// it has connection to another var of calls a function.

		if blankName(name.Name) {
			// TODO: Could a black name assignment have a side effect?
			//       Maybe if metrics aren't nil, give it a non-blank name.
			//		 var _ = func() bool { /*bad init*/ }()
			continue
		}

		tv, has := ab.curPkg.Source().TypesInfo.Defs[name]
		if !has {
			panic(terror.New(`value specification not found in types info`).
				With(`pos`, ab.pos(spec.Pos())))
		}

		typ := ab.convertType(tv.Type())
		ab.proj.NewValueDecl(constructs.ValueDeclArgs{
			Package:  ab.curPkg,
			Name:     name.Name,
			Const:    isConst,
			Type:     typ,
			Location: ab.proj.NewLoc(spec.Pos()),
		})
	}
}

func (ab *abstractor) setTypeParamOverrides(args *types.TypeList, params *types.TypeParamList, decl *ast.FuncDecl) {
	count := args.Len()
	if count != params.Len() {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, ab.pos(decl.Pos())))
	}

	ab.typeParamReplacer = map[*types.TypeParam]*types.TypeParam{}
	for i := range count {
		ab.typeParamReplacer[args.At(i).(*types.TypeParam)] = params.At(i)
	}
}

func (ab *abstractor) clearTypeParamOverrides() {
	ab.typeParamReplacer = nil
}

func (ab *abstractor) abstractReceiver(decl *ast.FuncDecl) (bool, string) {
	if decl.Recv == nil || decl.Recv.NumFields() <= 0 {
		return false, ``
	}

	if decl.Recv.NumFields() != 1 {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, ab.pos(decl.Pos())))
	}

	noCopyRecv := false
	tv, has := ab.curPkg.Source().TypesInfo.Types[decl.Recv.List[0].Type]
	if !has {
		panic(terror.New(`function receiver not found in types info`).
			With(`pos`, ab.pos(decl.Pos())))
	}

	recv := tv.Type
	if p, ok := recv.(*types.Pointer); ok {
		noCopyRecv = true
		recv = p.Elem()
	}

	n, ok := recv.(*types.Named)
	if !ok {
		panic(terror.New(`function declaration has unexpected receiver type`).
			WithType(`receiver`, recv).
			With(`pos`, ab.pos(decl.Pos())))
	}
	ab.setTypeParamOverrides(n.TypeArgs(), n.TypeParams(), decl)

	recvName := n.Origin().Obj().Name()
	return noCopyRecv, recvName
}

func (ab *abstractor) abstractFuncDecl(decl *ast.FuncDecl) {
	info := ab.curPkg.Source().TypesInfo
	obj := info.Defs[decl.Name]

	noCopyRecv, recvName := ab.abstractReceiver(decl)
	sig := ab.convertSignature(obj.Type().(*types.Signature))
	ab.clearTypeParamOverrides()

	mets := metrics.New(ab.curPkg.Source().Fset, decl)
	loc := ab.proj.NewLoc(decl.Pos())

	name := decl.Name.Name
	if name == `init` && len(recvName) <= 0 && sig.Vacant() {
		name = `init#` + strconv.Itoa(ab.curPkg.InitCount())
	}

	ab.proj.NewMethod(constructs.MethodArgs{
		Package:    ab.curPkg,
		Name:       name,
		Signature:  sig,
		Metrics:    mets,
		NoCopyRecv: noCopyRecv,
		RecvName:   recvName,
		Location:   loc,
	})
}
