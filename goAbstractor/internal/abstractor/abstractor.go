package abstractor

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"
	"reflect"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/wrapKind"
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

func convertList[T, U any](n int, getter func(i int) T, convert func(value T) *U) []*U {
	list := make([]*U, 0, n)
	for i := range n {
		if p := convert(getter(i)); p != nil {
			list = append(list, p)
		}
	}
	return slices.Compact(list)
}

func (ab *abstractor) convertType(t types.Type) typeDesc.TypeDesc {
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
		return ab.convertSignature(t2)
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

func (ab *abstractor) convertArray(t *types.Array) *typeDesc.Wrap {
	return &typeDesc.Wrap{
		Kind: wrapKind.Array,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertBasic(t *types.Basic) *typeDesc.Ref {
	return &typeDesc.Ref{
		Ref: t.Name(),
	}
}

func (ab *abstractor) convertChan(t *types.Chan) *typeDesc.Wrap {
	return &typeDesc.Wrap{
		Kind: wrapKind.Chan,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertInterface(t *types.Interface) *typeDesc.Interface {
	t = t.Complete()
	return ab.registerInterface(&typeDesc.Interface{
		Methods: convertList(t.NumMethods(), t.Method, ab.convertFunc),
	})
}

func (ab *abstractor) convertMap(t *types.Map) *typeDesc.Map {
	return &typeDesc.Map{
		Key:   ab.convertType(t.Key()),
		Value: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertNamed(t *types.Named) *typeDesc.Ref {
	return &typeDesc.Ref{
		Ref: t.String(),
	}
}

func (ab *abstractor) convertFunc(t *types.Func) *typeDesc.Func {
	return &typeDesc.Func{
		Name:      t.Name(),
		Signature: ab.convertSignature(t.Type().(*types.Signature)),
	}
}

func (ab *abstractor) convertPointer(t *types.Pointer) *typeDesc.Wrap {
	return &typeDesc.Wrap{
		Kind: wrapKind.Pointer,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertSignature(t *types.Signature) *typeDesc.Signature {
	// Don't output receiver or receiver type here.
	return ab.registerSignature(&typeDesc.Signature{
		Variadic:   t.Variadic(),
		Params:     ab.convertParamTuple(t.Params()),
		Return:     ab.createReturn(ab.convertFieldTuple(t.Results())),
		TypeParams: ab.convertTypeParamList(t.TypeParams()),
	})
}

func (ab *abstractor) convertSlice(t *types.Slice) *typeDesc.Wrap {
	return &typeDesc.Wrap{
		Kind: wrapKind.List,
		Elem: ab.convertType(t.Elem()),
	}
}

func (ab *abstractor) convertStruct(t *types.Struct) *typeDesc.Struct {
	return ab.registerStruct(&typeDesc.Struct{
		Fields: convertList(t.NumFields(), t.Field, ab.convertField),
	})
}

func uniqueName(names collections.Set[string]) string {
	for offset := 1; offset < 10000; offset++ {
		name := fmt.Sprintf(`value%d`, offset)
		if !names.Contains(name) {
			names.Add(name)
			return name
		}
	}
	return `_`
}

func (ab *abstractor) createReturn(returns []*typeDesc.Field) typeDesc.TypeDesc {
	switch len(returns) {
	case 0:
		return nil
	case 1:
		return returns[0].Type
	default:
		names := set.From(enumerator.Select(enumerator.Enumerate(returns...),
			func(f *typeDesc.Field) string { return f.Name }).NotZero())
		for _, f := range returns {
			f.Anonymous = false
			if len(f.Name) <= 0 || f.Name == `_` || f.Name == `.` {
				f.Name = uniqueName(names)
			}
		}
		return ab.registerStruct(&typeDesc.Struct{
			Fields: returns,
		})
	}
}

func (ab *abstractor) convertParamTuple(t *types.Tuple) []*typeDesc.Param {
	return convertList(t.Len(), t.At, ab.convertParam)
}

func (ab *abstractor) convertFieldTuple(t *types.Tuple) []*typeDesc.Field {
	return convertList(t.Len(), t.At, ab.convertField)
}

func (ab *abstractor) convertTypeParam(t *types.TypeParam) *typeDesc.TypeParam {
	return ab.registerTypeParam(&typeDesc.TypeParam{
		Index:      t.Index(),
		Constraint: ab.convertType(t.Constraint()),
	})
}

func (ab *abstractor) convertTypeParamList(t *types.TypeParamList) []*typeDesc.TypeParam {
	return convertList(t.Len(), t.At, ab.convertTypeParam)
}

func (ab *abstractor) convertParam(t *types.Var) *typeDesc.Param {
	return &typeDesc.Param{
		Name: t.Name(),
		Type: ab.convertType(t.Type()),
	}
}

func (ab *abstractor) convertField(t *types.Var) *typeDesc.Field {
	return &typeDesc.Field{
		Anonymous: t.Anonymous(),
		Name:      t.Name(),
		Type:      ab.convertType(t.Type()),
	}
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
