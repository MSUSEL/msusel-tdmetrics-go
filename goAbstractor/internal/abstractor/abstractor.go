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
)

// TODO:
// - Figure out implemented interfaces.
// - Determine what to do with pointer receivers to make it similar o Java.
// - Add analytics:
//   - Add cyclomatic complexity per method.
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)

func Abstract(ps []*packages.Package) *constructs.Project {
	ab := &abstractor{}
	return ab.abstractProject(ps)
}

type abstractor struct {
	allStructs    []*typeDesc.Struct
	allInterfaces []*typeDesc.Interface
	allSignatures []*typeDesc.Signature
	allTypeParam  []*typeDesc.TypeParam
}

func (ab *abstractor) abstractProject(ps []*packages.Package) *constructs.Project {
	proj := &constructs.Project{}
	packages.Visit(ps, func(src *packages.Package) bool {
		if ap := ab.abstractPackage(src); ap != nil {
			proj.Packages = append(proj.Packages, ap)
		}
		return true
	}, nil)

	// TODO: Add allStructs, allInterfaces, allTypeParam, and allSignatures

	return proj
}

func (ab *abstractor) abstractPackage(src *packages.Package) *constructs.Package {
	pkg := &constructs.Package{
		Path:    src.PkgPath,
		Imports: utils.SortedKeys(src.Imports),
	}
	for _, f := range src.Syntax {
		ab.addFile(pkg, src, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg *constructs.Package, src *packages.Package, f *ast.File) {
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
	for _, spec := range decl.Specs {
		switch s := spec.(type) {
		case *ast.ImportSpec:
			// ignore
		case *ast.TypeSpec:
			ab.abstractTypeSpec(pkg, src, s)
		case *ast.ValueSpec:
			ab.abstractValueSpec(pkg, src, s)
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

func (ab *abstractor) abstractValueSpec(pkg *constructs.Package, src *packages.Package, spec *ast.ValueSpec) {
	for _, name := range spec.Names {
		if name.Name != `_` {
			tv, has := src.TypesInfo.Defs[name]
			if !has {
				panic(fmt.Errorf(`value specification not found in types info: %s`, pos(src, spec.Type.Pos())))
			}

			def := &constructs.ValueDef{
				Name: name.Name,
				Type: ab.convertType(tv.Type()),
			}
			pkg.Values = append(pkg.Values, def)
		}
	}
}

func (ab *abstractor) abstractFuncDecl(pkg *constructs.Package, src *packages.Package, decl *ast.FuncDecl) {
	obj := src.TypesInfo.Defs[decl.Name]
	m := &constructs.Method{
		Name:      decl.Name.Name,
		Signature: ab.convertSignature(obj.Type().(*types.Signature)),
	}

	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(fmt.Errorf(`function declaration has unexpected receiver fields: %s`, pos(src, decl.Pos())))
		}
		tv := src.TypesInfo.Types[decl.Recv.List[0].Type]
		m.Receiver = ab.convertType(tv.Type)
	}

	pkg.Methods = append(pkg.Methods, m)
}

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) registerStruct(s *typeDesc.Struct) *typeDesc.Struct {
	for _, s2 := range ab.allStructs {
		if reflect.DeepEqual(s, s2) {
			return s2
		}
	}
	ab.allStructs = append(ab.allStructs, s)
	return s
}

func (ab *abstractor) registerInterface(ti *typeDesc.Interface) *typeDesc.Interface {
	for _, t2 := range ab.allInterfaces {
		if reflect.DeepEqual(ti, t2) {
			return t2
		}
	}
	ab.allInterfaces = append(ab.allInterfaces, ti)
	return ti
}

func (ab *abstractor) registerSignature(sig *typeDesc.Signature) *typeDesc.Signature {
	for _, s2 := range ab.allSignatures {
		if reflect.DeepEqual(sig, s2) {
			return s2
		}
	}
	ab.allSignatures = append(ab.allSignatures, sig)
	return sig
}

func (ab *abstractor) registerTypeParam(tp *typeDesc.TypeParam) *typeDesc.TypeParam {
	for _, t2 := range ab.allTypeParam {
		if reflect.DeepEqual(tp, t2) {
			return t2
		}
	}
	ab.allTypeParam = append(ab.allTypeParam, tp)
	return tp
}
