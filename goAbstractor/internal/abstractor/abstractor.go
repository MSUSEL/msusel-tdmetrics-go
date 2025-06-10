package abstractor

import (
	"fmt"
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
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/querier"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/resolver"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/innate"
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
		log     = cfg.Log
		querier = querier.New(cfg.Packages)
		locs    = locs.NewSet(querier.FileSet())
		proj    = project.New(locs)
		bk      = baker.New(proj)
	)

	ab := &abstractor{
		querier:   querier,
		log:       log,
		baker:     bk,
		proj:      proj,
		typeCache: map[any]any{},
	}
	ab.abstractProject()

	resolver.Resolve(proj, log)

	log.Log(`done`)
	return proj
}

type abstractor struct {
	querier       *querier.Querier
	log           *logger.Logger
	baker         baker.Baker
	proj          constructs.Project
	curPkg        constructs.Package
	curNest       constructs.NestType
	implicitTypes []constructs.TypeDesc
	tpReplacer    map[*types.TypeParam]*types.TypeParam
	typeCache     map[any]any
}

func (ab *abstractor) pos(pos token.Pos) token.Position {
	return ab.querier.Pos(pos)
}

func (ab *abstractor) converter() converter.Converter {
	return converter.New(ab.log, ab.querier, ab.baker, ab.proj,
		ab.curPkg, ab.curNest, ab.implicitTypes, ab.tpReplacer, ab.typeCache)
}

func (ab *abstractor) abstractProject() {
	ab.log.Log(`abstract project`)
	log2 := ab.log.Group(`packages`).Prefix(`|  `)
	ab.querier.ForeachPackage(func(src *packages.Package) {
		ab.abstractPackage(src, log2)
	})
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
	t := ab.querier.GetType(spec.Type)
	context := t.String()
	loc := ab.proj.Locs().NewLoc(spec.Pos())
	tp := ab.abstractTypeParams(spec.TypeParams, context)
	typ := ab.converter().ConvertType(t, context)

	if it, ok := typ.(constructs.InterfaceDesc); ok {
		ab.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			RealType:   t,
			Package:    ab.curPkg,
			Name:       spec.Name.Name,
			Exported:   spec.Name.IsExported(),
			Interface:  it,
			TypeParams: tp,
			Location:   loc,
			Nest:       ab.curNest,
		})
		return
	}

	st, ok := typ.(constructs.StructDesc)
	if !ok {
		st = ab.proj.NewStructDesc(constructs.StructDescArgs{
			Fields: []constructs.Field{
				ab.proj.NewField(constructs.FieldArgs{
					Name:     innate.Data,
					Exported: true,
					Embedded: true,
					Type:     typ,
				}),
			},
			Package: ab.curPkg.Source(),
		})
	}

	ab.proj.NewObject(constructs.ObjectArgs{
		RealType:   t,
		Package:    ab.curPkg,
		Name:       spec.Name.Name,
		Exported:   spec.Name.IsExported(),
		Data:       st,
		TypeParams: tp,
		Location:   loc,
		Nest:       ab.curNest,
	})
}

func (ab *abstractor) abstractTypeParams(fields *ast.FieldList, context string) []constructs.TypeParam {
	ns := []constructs.TypeParam{}
	if !utils.IsNil(fields) {
		for _, field := range fields.List {
			ns = append(ns, ab.abstractTypeParam(field, context)...)
		}
	}
	return ns
}

func (ab *abstractor) abstractTypeParam(field *ast.Field, context string) []constructs.TypeParam {
	ns := []constructs.TypeParam{}
	if utils.IsNil(field) {
		return ns
	}

	t := ab.querier.GetType(field.Type)
	typ := ab.converter().ConvertType(t, context)
	for _, name := range field.Names {
		named := ab.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: name.Name,
			Type: typ,
		})
		ns = append(ns, named)
	}
	return ns
}

func (ab *abstractor) analyze(node ast.Node) constructs.Metrics {
	return analyzer.Analyze(ab.log, ab.querier, ab.proj, ab.curPkg, ab.baker, ab.converter(), node)
}

func (ab *abstractor) abstractValueSpec(spec *ast.ValueSpec, isConst bool) {
	var metrics constructs.Metrics
	for i, name := range spec.Names {
		if i < len(spec.Values) {
			metrics = ab.analyze(spec.Values[i])
		}

		obj := ab.querier.GetDef(name)
		typ := ab.converter().ConvertType(obj.Type(), name.Name)
		ab.proj.NewValue(constructs.ValueArgs{
			Package:  ab.curPkg,
			Name:     name.Name,
			Exported: name.IsExported(),
			Const:    isConst,
			Metrics:  metrics,
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

	ptrRecv := false
	recv := ab.querier.GetType(decl.Recv.List[0].Type)
	if p, ok := recv.(*types.Pointer); ok {
		ptrRecv = true
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
	return ptrRecv, recvName
}

func (ab *abstractor) abstractFuncDecl(decl *ast.FuncDecl) {
	info := ab.querier.Info()
	obj := info.Defs[decl.Name]
	loc := ab.proj.Locs().NewLoc(decl.Pos())

	ptrRecv, recvName := ab.abstractReceiver(decl)
	sig := ab.converter().ConvertSignature(obj.Type().(*types.Signature), decl.Name.Name)

	prevNest := ab.curNest
	defer func() { ab.curNest = prevNest }()
	ab.curNest = ab.proj.NewTempDeclRef(constructs.TempDeclRefArgs{
		PackagePath:   ab.curPkg.Path(),
		Name:          decl.Name.Name,
		ImplicitTypes: ab.implicitTypes,
	})

	fmt.Printf("-----------------------------\n")                                        // TODO: REMOVE
	fmt.Printf(">> abstractor: abstractFuncDecl: decl.Name:     %s\n", decl.Name.Name)   // TODO: REMOVE
	fmt.Printf(">> abstractor: abstractFuncDecl: nest:          %v\n", ab.curNest)       // TODO: REMOVE
	fmt.Printf(">> abstractor: abstractFuncDecl: implicitTypes: %v\n", ab.implicitTypes) // TODO: REMOVE

	metrics := ab.analyze(decl)
	ab.clearTypeParamOverrides()

	tp := ab.abstractTypeParams(decl.Type.TypeParams, decl.Name.Name)

	exported := decl.Name.IsExported()
	name := decl.Name.Name
	if name == `init` && len(recvName) <= 0 && sig.IsVacant() {
		name = `init#` + strconv.Itoa(ab.curPkg.InitCount())
	}

	method := ab.proj.NewMethod(constructs.MethodArgs{
		Package:     ab.curPkg,
		Name:        name,
		Exported:    exported,
		Location:    loc,
		TypeParams:  tp,
		Signature:   sig,
		Metrics:     metrics,
		RecvName:    recvName,
		PointerRecv: ptrRecv,
	})

	ab.curNest = method
	defer func() { ab.curNest = nil }()
	ab.abstractNestedTypes(decl.Body)
}

func (ab *abstractor) abstractNestedTypes(body *ast.BlockStmt) {
	if body == nil {
		return
	}
	// TODO: Simplify the following to only look for type specs.
	ast.Inspect(body, func(node ast.Node) bool {
		if stmt, ok := node.(*ast.DeclStmt); ok {
			if decl, ok := stmt.Decl.(*ast.GenDecl); ok {
				for _, spec := range decl.Specs {
					if typeSpec, ok := spec.(*ast.TypeSpec); ok {
						ab.abstractTypeSpec(typeSpec)
					}
				}
			}
		}
		return true
	})
}
