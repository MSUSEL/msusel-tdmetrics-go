package abstractor

import (
	"fmt"
	"go/types"
	"slices"
	"sort"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/set"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/wrapKind"
)

func convertList[T, U any](n int, getter func(i int) T, convert func(value T) *U) []*U {
	list := make([]*U, 0, n)
	for i := range n {
		if p := convert(getter(i)); p != nil {
			list = append(list, p)
		}
	}
	if len(list) <= 0 {
		return nil
	}
	return slices.Clip(list)
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
	methods := convertList(t.NumMethods(), t.Method, ab.convertFunc)
	// Sort interface methods since order doesn't matter.
	sort.SliceIsSorted(methods, func(i, j int) bool {
		return methods[i].Name < methods[j].Name
	})
	return ab.registerInterface(&typeDesc.Interface{
		Methods: methods,
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
			if len(f.Name) <= 0 || f.Name == `_` {
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
