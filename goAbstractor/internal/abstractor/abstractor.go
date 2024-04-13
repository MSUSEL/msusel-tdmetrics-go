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
// - Determine what to do with pointer receivers to make it similar to Java.
// - Add analytics:
//   - Add cyclomatic complexity per method.
//   - The set of variables with locations that are read from and written
//     to in each method. Used in Tight Class Cohesion (TCC) and
//     Design Recovery (DR).
//   - The set of all methods called in each method. Used for
//     Access to Foreign Data (ATFD) and Design Recovery (DR)

func Abstract(ps []*packages.Package, verbose bool) *constructs.Project {
	ab := &abstractor{
		verbose: verbose,
		proj:    &constructs.Project{},
	}
	ab.abstractProject(ps)
	ab.resolveExtends()
	ab.resolveImplementation()

	// Leave indices as zero until the end so that checking equality
	// using reflect.DeepEqual will not find differences in the indices.
	ab.updateIndices()
	return ab.proj
}

type abstractor struct {
	verbose bool
	proj    *constructs.Project
}

func (ab *abstractor) abstractProject(ps []*packages.Package) {
	packages.Visit(ps, func(src *packages.Package) bool {
		if ap := ab.abstractPackage(src); ap != nil {
			ab.proj.Packages = append(ab.proj.Packages, ap)
		}
		return true
	}, nil)
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
	}

	if decl.Recv != nil && decl.Recv.NumFields() > 0 {
		if decl.Recv.NumFields() != 1 {
			panic(fmt.Errorf(`function declaration has unexpected receiver fields: %s`, pos(src, decl.Pos())))
		}
		recv := src.TypesInfo.Types[decl.Recv.List[0].Type].Type
		// Ignore the pointer since abstraction doesn't need
		// to know if a reference or pointer is being set.
		if p, ok := recv.(*types.Pointer); ok {
			recv = p.Elem()
		}
		m.Receiver = ab.convertType(recv)
	}

	// TODO: Evaluate the body of the function to abstract metrics.
	pkg.Methods = append(pkg.Methods, m)
}

func pos(src *packages.Package, pos token.Pos) string {
	return src.Fset.Position(pos).String()
}

func (ab *abstractor) registerInterface(ti *typeDesc.Interface) *typeDesc.Interface {
	for _, t2 := range ab.proj.AllInterfaces {
		if reflect.DeepEqual(ti, t2) {
			return t2
		}
	}
	ab.proj.AllInterfaces = append(ab.proj.AllInterfaces, ti)
	return ti
}

func (ab *abstractor) registerSignature(sig *typeDesc.Signature) *typeDesc.Signature {
	for _, s2 := range ab.proj.AllSignatures {
		if reflect.DeepEqual(sig, s2) {
			return s2
		}
	}
	ab.proj.AllSignatures = append(ab.proj.AllSignatures, sig)
	return sig
}

func (ab *abstractor) registerStruct(s *typeDesc.Struct) *typeDesc.Struct {
	for _, s2 := range ab.proj.AllStructs {
		if reflect.DeepEqual(s, s2) {
			return s2
		}
	}
	ab.proj.AllStructs = append(ab.proj.AllStructs, s)
	return s
}

func (ab *abstractor) registerTypeParam(tp *typeDesc.TypeParam) *typeDesc.TypeParam {
	for _, t2 := range ab.proj.AllTypeParams {
		if reflect.DeepEqual(tp, t2) {
			return t2
		}
	}
	ab.proj.AllTypeParams = append(ab.proj.AllTypeParams, tp)
	return tp
}

func (ab *abstractor) resolveExtends() {
	// TODO: Finish
}

func (ab *abstractor) resolveImplementation() {
	// TODO: Finish
}

func (ab *abstractor) updateIndices() {
	offset := 0
	for i, c := range ab.proj.AllInterfaces {
		c.Index = i + offset
	}
	offset += len(ab.proj.AllInterfaces)
	for i, c := range ab.proj.AllSignatures {
		c.Index = i + offset
	}
	offset += len(ab.proj.AllSignatures)
	for i, c := range ab.proj.AllStructs {
		c.Index = i + offset
	}
	offset += len(ab.proj.AllStructs)
	for i, c := range ab.proj.AllTypeParams {
		c.Index = i + offset
	}
}
