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

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/analyzer"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/baker"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/converter"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/project"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/logger"
)

type Config struct {
	Packages []*packages.Package
	Log      *logger.Logger
}

func Abstract(cfg Config) constructs.Project {
	var (
		log  = cfg.Log
		fSet = cfg.Packages[0].Fset
		locs = locs.NewSet(fSet)
		proj = project.New(locs)
		bk   = baker.New(proj)
	)

	ab := &abstractor{
		packages: cfg.Packages,
		log:      log,
		baker:    bk,
		proj:     proj,
	}
	ab.abstractProject()

	resolver.Resolve(proj, log)

	log.Log(`done`)
	return proj
}

type abstractor struct {
	packages   []*packages.Package
	log        *logger.Logger
	baker      baker.Baker
	proj       constructs.Project
	curPkg     constructs.Package
	tpReplacer map[*types.TypeParam]*types.TypeParam
}

func (ab *abstractor) pos(pos token.Pos) token.Position {
	return ab.curPkg.Source().Fset.Position(pos)
}

func (ab *abstractor) info() *types.Info {
	return ab.curPkg.Source().TypesInfo
}

func (ab *abstractor) converter() converter.Converter {
	return converter.New(ab.baker, ab.proj, ab.curPkg, ab.tpReplacer)
}

func (ab *abstractor) abstractProject() {
	ab.log.Log(`abstract project`)
	log2 := ab.log.Group(`packages`).Prefix(`|  `)
	packages.Visit(ab.packages, func(src *packages.Package) bool {
		ab.abstractPackage(src, log2)
		return true
	}, nil)
}

func (ab *abstractor) abstractPackage(src *packages.Package, log *logger.Logger) {
	log.Logf(`abstract package: %s`, src.PkgPath)
	ab.curPkg = ab.proj.NewPackage(constructs.PackageArgs{
		RealPkg:     src,
		Path:        src.PkgPath,
		Name:        src.Name,
		ImportPaths: utils.SortedKeys(src.Imports),
	})
	log2 := log.Group(`files`).Prefix(`|  `)
	for _, f := range src.Syntax {
		ab.abstractFile(f, log2)
	}
}

func (ab *abstractor) abstractFile(f *ast.File, log *logger.Logger) {
	path := ab.pos(f.FileStart).Filename
	basePath := filepath.Base(path)
	if ab.curPkg.EntryPoint() {
		ab.proj.Locs().Alias(path, basePath)
	} else {
		pkgPath := ab.curPkg.Source().PkgPath
		alias := filepath.ToSlash(filepath.Join(pkgPath, basePath))
		ab.proj.Locs().Alias(path, alias)
	}

	log.Logf(`add file to package: %s`, basePath)
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
	tv, has := ab.info().Types[spec.Type]
	if !has {
		panic(terror.New(`type specification not found in types info`).
			With(`pos`, ab.pos(spec.Pos())))
	}

	loc := ab.proj.Locs().NewLoc(spec.Pos())
	tp := ab.abstractTypeParams(spec.TypeParams)
	typ := ab.converter().ConvertType(tv.Type)

	if it, ok := typ.(constructs.InterfaceDesc); ok {
		ab.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			RealType:   tv.Type,
			Package:    ab.curPkg,
			Name:       spec.Name.Name,
			Exported:   spec.Name.IsExported(),
			Interface:  it,
			TypeParams: tp,
			Location:   loc,
		})
		return
	}

	st, ok := typ.(constructs.StructDesc)
	if !ok {
		st = ab.proj.NewStructDesc(constructs.StructDescArgs{
			Fields: []constructs.Field{
				ab.proj.NewField(constructs.FieldArgs{
					Name:     `$data`,
					Exported: true,
					Embedded: true,
					Type:     typ,
				}),
			},
			Package: ab.curPkg.Source(),
		})
	}

	ab.proj.NewObject(constructs.ObjectArgs{
		RealType:   tv.Type,
		Package:    ab.curPkg,
		Name:       spec.Name.Name,
		Exported:   spec.Name.IsExported(),
		Data:       st,
		TypeParams: tp,
		Location:   loc,
	})
}

func (ab *abstractor) abstractTypeParams(fields *ast.FieldList) []constructs.TypeParam {
	ns := []constructs.TypeParam{}
	if !utils.IsNil(fields) {
		for _, field := range fields.List {
			ns = append(ns, ab.abstractTypeParam(field)...)
		}
	}
	return ns
}

func (ab *abstractor) abstractTypeParam(field *ast.Field) []constructs.TypeParam {
	ns := []constructs.TypeParam{}
	if utils.IsNil(field) {
		return ns
	}

	tv, has := ab.info().Types[field.Type]
	if !has {
		panic(terror.New(`field not found in types info`).
			With(`pos`, ab.pos(field.Pos())))
	}

	typ := ab.converter().ConvertType(tv.Type)
	for _, name := range field.Names {
		named := ab.proj.NewTypeParam(constructs.TypeParamArgs{
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

		if constructs.BlankName(name.Name) {
			// TODO: Could a black name assignment have a side effect?
			//       Maybe if metrics aren't nil, give it a non-blank name.
			//		 var _ = func() bool { /*pseudo init*/ }()
			continue
		}

		tv, has := ab.info().Defs[name]
		if !has {
			panic(terror.New(`value specification not found in types info`).
				With(`pos`, ab.pos(spec.Pos())))
		}

		typ := ab.converter().ConvertType(tv.Type())
		ab.proj.NewValue(constructs.ValueArgs{
			Package:  ab.curPkg,
			Name:     name.Name,
			Exported: name.IsExported(),
			Const:    isConst,
			Type:     typ,
			Location: ab.proj.Locs().NewLoc(spec.Pos()),
		})
	}
}

func (ab *abstractor) setTypeParamOverrides(args *types.TypeList, params *types.TypeParamList, decl *ast.FuncDecl) {
	count := args.Len()
	if count != params.Len() {
		panic(terror.New(`function declaration has unexpected receiver fields`).
			With(`pos`, ab.pos(decl.Pos())))
	}

	ab.tpReplacer = map[*types.TypeParam]*types.TypeParam{}
	for i := range count {
		tp := args.At(i).(*types.TypeParam)
		ab.tpReplacer[tp] = params.At(i)
	}
}

func (ab *abstractor) clearTypeParamOverrides() {
	ab.tpReplacer = nil
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
	tv, has := ab.info().Types[decl.Recv.List[0].Type]
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
	info := ab.info()
	obj := info.Defs[decl.Name]

	noCopyRecv, recvName := ab.abstractReceiver(decl)
	sig := ab.converter().ConvertSignature(obj.Type().(*types.Signature))
	ab.clearTypeParamOverrides()

	an := analyzer.New(ab.proj.Locs()).Analyze(decl)
	metrics := ab.proj.NewMetrics(an.GetMetrics())
	loc := ab.proj.Locs().NewLoc(decl.Pos())
	tp := ab.abstractTypeParams(decl.Type.TypeParams)

	name := decl.Name.Name
	if name == `init` && len(recvName) <= 0 && sig.IsVacant() {
		name = `init#` + strconv.Itoa(ab.curPkg.InitCount())
	}

	ab.proj.NewMethod(constructs.MethodArgs{
		Package:    ab.curPkg,
		Name:       name,
		Exported:   decl.Name.IsExported(),
		Location:   loc,
		TypeParams: tp,
		Signature:  sig,
		Metrics:    metrics,
		RecvName:   recvName,
		NoCopyRecv: noCopyRecv,
	})
}
