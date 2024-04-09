package abstractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/construct"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/construct/typeKind"
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

func Abstract(ps []*packages.Package) *construct.Project {
	ab := &abstractor{}
	return ab.abstractProject(ps)
}

type abstractor struct {
	allStructs    []*construct.Struct
	allInterfaces []*construct.Interface
	allSignatures []*construct.Signature
	allTypeParam  []*construct.TypeParam
}

func (ab *abstractor) abstractProject(ps []*packages.Package) *construct.Project {
	proj := &construct.Project{}
	packages.Visit(ps, func(src *packages.Package) bool {
		if ap := ab.abstractPackage(src); ap != nil {
			proj.Packages = append(proj.Packages, ap)
		}
		return true
	}, nil)

	// TODO: Add allStructs, allInterfaces, allTypeParam, and allSignatures

	return proj
}

func (ab *abstractor) abstractPackage(src *packages.Package) *construct.Package {
	pkg := &construct.Package{
		Path:    src.PkgPath,
		Imports: utils.SortedKeys(src.Imports),
	}
	for _, f := range src.Syntax {
		ab.addFile(pkg, src, f)
	}
	return pkg
}

func (ab *abstractor) addFile(pkg *construct.Package, src *packages.Package, f *ast.File) {
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

func (ab *abstractor) addGenDecl(pkg *construct.Package, src *packages.Package, decl *ast.GenDecl) {
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

func (ab *abstractor) abstractTypeSpec(pkg *construct.Package, src *packages.Package, spec *ast.TypeSpec) {
	tv, has := src.TypesInfo.Types[spec.Type]
	if !has {
		panic(fmt.Errorf(`type specification not found in types info: %s`, pos(src, spec.Type.Pos())))
	}
	def := &construct.TypeDef{
		Name: spec.Name.Name,
		Type: ab.convertType(tv.Type),
	}
	pkg.Types = append(pkg.Types, def)
}

func (ab *abstractor) abstractValueSpec(pkg *construct.Package, src *packages.Package, spec *ast.ValueSpec) {
	for _, name := range spec.Names {
		if name.Name != `_` {
			tv, has := src.TypesInfo.Defs[name]
			if !has {
				panic(fmt.Errorf(`value specification not found in types info: %s`, pos(src, spec.Type.Pos())))
			}

			def := &construct.ValueDef{
				Name: name.Name,
				Type: ab.convertType(tv.Type()),
			}
			pkg.Values = append(pkg.Values, def)
		}
	}
}

func (ab *abstractor) abstractFuncDecl(pkg *construct.Package, src *packages.Package, decl *ast.FuncDecl) {
	obj := src.TypesInfo.Defs[decl.Name]
	m := &construct.Method{
		Name:      decl.Name.Name,
		Signature: ab.convertSignature(obj.Type().(*types.Signature), false),
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

func convertList[T, U any](n int, getter func(i int) T, convert func(value T) *U) []*U {
	list := make([]*U, 0, n)
	for i := range n {
		if p := convert(getter(i)); p != nil {
			list = append(list, p)
		}
	}
	return slices.Compact(list)
}

func (ab *abstractor) convertType(t types.Type) construct.TypeDesc {
	switch t2 := t.(type) {
	case *types.Array:
		return ab.convertArray(t2)
	case *types.Basic:
		return ab.convertBasic(t2)
	case *types.Chan:
		return ab.convertChan(t2)
	case *types.Interface:
		return ab.convertInterface(t2)
	case *types.Map:
		return ab.convertMap(t2)
	case *types.Named:
		return ab.convertNamed(t2)
	case *types.Pointer:
		return ab.convertPointer(t2)
	case *types.Signature:
		return ab.convertSignature(t2, true)
	case *types.Slice:
		return ab.convertSlice(t2)
	case *types.Struct:
		return ab.convertStruct(t2)
	case *types.TypeParam:
		return ab.convertTypeParam(t2)
	default:
		panic(fmt.Errorf(`unhandled type, %T: %s`, t, t))
	}
}

func (ab *abstractor) convertArray(t *types.Array) *construct.TypeWrap {
	return &construct.TypeWrap{
		Kind: typeKind.Array,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertBasic(t *types.Basic) *construct.TypeRef {
	return &construct.TypeRef{
		Ref: t.Name(),
	}
}

func (ab *abstractor) convertChan(t *types.Chan) *construct.TypeWrap {
	return &construct.TypeWrap{
		Kind: typeKind.Chan,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertInterface(t *types.Interface) *construct.Interface {
	t = t.Complete()
	return ab.registerInterface(&construct.Interface{
		Methods: convertList(t.NumMethods(), t.Method, ab.convertFunc),
	})
}

func (ab *abstractor) convertMap(t *types.Map) *construct.TypeMap {
	return &construct.TypeMap{
		Key:   ab.convertType(t.Key()),
		Value: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertNamed(t *types.Named) *construct.TypeRef {
	return &construct.TypeRef{
		Ref: t.String(),
	}
}

func (ab *abstractor) convertFunc(t *types.Func) *construct.TypeFunc {
	return &construct.TypeFunc{
		Name:      t.Name(),
		Signature: ab.convertSignature(t.Type().(*types.Signature), false),
	}
}

func (ab *abstractor) convertPointer(t *types.Pointer) *construct.TypeWrap {
	return &construct.TypeWrap{
		Kind: typeKind.Pointer,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertSignature(t *types.Signature, showKind bool) *construct.Signature {
	// Don't output receiver or receiver type here.
	return ab.registerSignature(&construct.Signature{
		ShowKind:   showKind,
		Variadic:   t.Variadic(),
		Params:     ab.convertTuple(t.Params()),
		Returns:    ab.convertTuple(t.Results()),
		TypeParams: ab.convertTypeParamList(t.TypeParams()),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) *construct.TypeWrap {
	return &construct.TypeWrap{
		Kind: typeKind.List,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertStruct(t *types.Struct) *construct.Struct {
	return ab.registerStruct(&construct.Struct{
		Fields: convertList(t.NumFields(), t.Field, ab.convertVar),
	})
}

func (ab *abstractor) convertTuple(t *types.Tuple) []*construct.TypeVar {
	return convertList(t.Len(), t.At, ab.convertVar)
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) *construct.TypeParam {
	return ab.registerTypeParam(&construct.TypeParam{
		Index:      t.Index(),
		Constraint: ab.convertType(t.Constraint()),
	})
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []*construct.TypeParam {
	return convertList(t.Len(), t.At, ab.convertTypeParam)
}

func (ab *abstractor) convertVar(t *types.Var) *construct.TypeVar {
	return &construct.TypeVar{
		Name: t.Name(),
		Type: ab.convertType(t.Type()),
	}
}

func (ab *abstractor) registerStruct(s *construct.Struct) *construct.Struct {
	// TODO: FINISH
	ab.allStructs = append(ab.allStructs, s)
	return s
}

func (ab *abstractor) registerInterface(ti *construct.Interface) *construct.Interface {
	// TODO: FINISH
	ab.allInterfaces = append(ab.allInterfaces, ti)
	return ti
}

func (ab *abstractor) registerSignature(sig *construct.Signature) *construct.Signature {
	// TODO: FINISH
	ab.allSignatures = append(ab.allSignatures, sig)
	return sig
}

func (ab *abstractor) registerTypeParam(tp *construct.TypeParam) *construct.TypeParam {
	// TODO: FINISH
	ab.allTypeParam = append(ab.allTypeParam, tp)
	return tp
}
