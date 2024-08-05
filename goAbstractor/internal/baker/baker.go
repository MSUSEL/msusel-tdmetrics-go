// The baked in types are stored for quick lookup and
// to ensure only one instance of each type is created.
//
// The function names are prepended with a `$` to avoid duck-typing
// with user-defined types. Some of these types don't normally exist in Go,
// but are used to represent the construct in a way that can be abstracted.
// Other types represent the built-in types such as error.
package baker

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

const BuiltinName = `$builtin`

type (
	Baker interface {
		BakeBuiltin() constructs.Package
		BakeBasic(typeName string) constructs.Basic
		BakeAny() constructs.Object
		BakeList() constructs.Object
		BakeChan() constructs.Object
		BakeMap() constructs.Object
		BakePointer() constructs.Object
		BakeComplex64() constructs.Object
		BakeComplex128() constructs.Object
		BakeError() constructs.Object
		BakeComparable() constructs.Object
	}

	bakerImp struct {
		fSet  *token.FileSet
		proj  constructs.Project
		baked map[string]any
	}
)

func New(fSet *token.FileSet, proj constructs.Project) Baker {
	return &bakerImp{
		fSet:  fSet,
		proj:  proj,
		baked: map[string]any{},
	}
}

func bakeOnce[T any](b *bakerImp, key string, create func() T) T {
	if baked, has := b.baked[key]; has {
		t, ok := baked.(T)
		if !ok {
			panic(terror.New(`unexpected baked type`).
				With(`key`, key).
				WithType(`wanted`, utils.Zero[T]()).
				WithType(`gotten type`, baked).
				With(`gotten value`, baked))
		}
		return t
	}

	t := create()
	b.baked[key] = t
	return t
}

// bakeIntFunc bakes in a signature for `func() int`.
// This is useful for things like `cap() int` and `len() int`.
func (b *bakerImp) bakeIntFunc() constructs.Signature {
	return bakeOnce(b, `func() int`, func() constructs.Signature {
		// <unnamed> int
		result := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`int`),
		})

		// func() int
		return b.proj.NewSignature(constructs.SignatureArgs{
			Package: b.BakeBuiltin(),
			Results: []constructs.Argument{result},
		})
	})
}

// BakeBuiltin bakes in a package to represent the builtin package.
func (b *bakerImp) BakeBuiltin() constructs.Package {
	return bakeOnce(b, BuiltinName, func() constructs.Package {
		builtinPkg := &packages.Package{
			PkgPath: BuiltinName,
			Name:    BuiltinName,
			Fset:    b.fSet,
			Types:   types.NewPackage(BuiltinName, BuiltinName),
		}

		return b.proj.NewPackage(constructs.PackageArgs{
			RealPkg: builtinPkg,
			Path:    BuiltinName,
			Name:    BuiltinName,
		})
	})
}

// BakeBasic bakes in a basic type by name. The name must
// be a valid basic type (e.g. int, int32, float64)
// but may not be complex numbers or interfaces like any.
func (b *bakerImp) BakeBasic(typeName string) constructs.Basic {
	return bakeOnce(b, `basic `+typeName, func() constructs.Basic {
		return b.proj.NewBasic(constructs.BasicArgs{
			Package:  b.BakeBuiltin(),
			TypeName: typeName,
		})
	})
}

// BakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (b *bakerImp) BakeAny() constructs.Object {
	return bakeOnce(b, `any`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		// any
		return b.proj.NewObject(constructs.ObjectArgs{
			RealType: types.NewInterfaceType(nil, nil),
			Package:  pkg,
			Name:     `any`,
			Location: locs.NoLoc(),
		})
	})
}

// bakeList bakes in an interface to represent a Go array or slice:
//
//	type list[T any] interface {
//		$len() int
//		$cap() int
//		$get(index int) T
//		$set(index int, value T)
//	}
//
// Note: The difference between an array and slice aren't
// important for the abstractor, so they are combined into one.
func (b *bakerImp) BakeList() constructs.Object {
	return bakeOnce(b, `List[T any]`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		tps := []constructs.TypeParam{tp}

		lenCapFunc := b.bakeIntFunc()
		indexParam := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `index`,
			Type: b.BakeBasic(`int`),
		})
		valueParam := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: tp,
		})
		valueResult := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: tp,
		})

		getFunc := b.proj.NewSignature(constructs.SignatureArgs{
			Params:  []constructs.Argument{indexParam},
			Results: []constructs.Argument{valueResult},
		})
		setFunc := b.proj.NewSignature(constructs.SignatureArgs{
			Params: []constructs.Argument{indexParam, valueParam},
		})

		// List[T]
		obj := b.proj.NewObject(constructs.ObjectArgs{
			Package:    pkg,
			Name:       `List`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
		})

		// $len() int
		b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       `$len`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Signature:  lenCapFunc,
			Receiver:   obj,
		})

		// $cap() int
		b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       `$cap`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Signature:  lenCapFunc,
			Receiver:   obj,
		})

		// $get(index int) T
		b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       `$get`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Signature:  getFunc,
			Receiver:   obj,
		})

		// $set(index int, value T)
		b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       `$set`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Signature:  setFunc,
			Receiver:   obj,
		})

		return obj
	})
}

// BakeChan bakes in an interface to represent a Go chan:
//
//	type chan[T any] interface {
//		$len() int
//		$recv() (T, bool)
//		$send(value T)
//	}
//
// If the given elements is nil, then the generic form is returned.
// Otherwise, the instance realization on the given element is returned.
//
// Note: Doesn't currently have cap, trySend, or tryRecv as defined in reflect.
func (b *bakerImp) BakeChan() constructs.Object {
	return bakeOnce(b, `Chan[T any]`, func() constructs.Object {
		pkg := b.BakeBuiltin()
		tp := b.proj.NewNamed(constructs.NamedArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = b.bakeIntFunc()

		// $recv() (T, bool)
		methods[`$recv`] = b.proj.NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     b.bakeReturnTuple(tp),
			Package:    pkg,
		})

		// value T
		valueParam := b.proj.NewNamed(constructs.NamedArgs{
			Name: `value`,
			Type: tp,
		})

		// $send(value T)
		methods[`$send`] = b.proj.NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{valueParam},
			Package:    pkg,
		})

		// chan[T any] interface
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    b.toNamedList(methods),
			Package:    pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `Chan`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakeMap bakes in an interface to represent a Go map:
//
//	type map[TKey, TValue any] interface {
//		$len() int
//		$get(key TKey) (TValue, bool)
//		$set(key TKey, value TValue)
//	}
//
// Note: Doesn't currently require Key to be comparable as defined in reflect.
func (b *bakerImp) BakeMap() constructs.Object {
	return bakeOnce(b, `Map[TKey comparable, TValue any]`, func() constructs.Object {
		pkg := b.BakeBuiltin()
		tpKey := b.proj.NewNamed(constructs.NamedArgs{
			Name: `TKey`,
			Type: b.BakeComparable(),
		})
		tpValue := b.proj.NewNamed(constructs.NamedArgs{
			Name: `TValue`,
			Type: b.BakeAny(),
		})
		tp := []constructs.Named{tpKey, tpValue}
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = b.bakeIntFunc()

		// key TKey
		keyParam := b.proj.NewNamed(constructs.NamedArgs{
			Name: `key`,
			Type: tpKey,
		})

		// $get(key TKey) (TValue, bool)
		methods[`$get`] = b.proj.NewSignature(constructs.SignatureArgs{
			TypeParams: tp,
			Params:     []constructs.Named{keyParam},
			Return:     b.bakeReturnTuple(tpValue),
			Package:    pkg,
		})

		// value TValue
		valueParam := b.proj.NewNamed(constructs.NamedArgs{
			Name: `value`,
			Type: tpValue,
		})

		// $set(key TKey, value TValue)
		methods[`$set`] = b.proj.NewSignature(constructs.SignatureArgs{
			TypeParams: tp,
			Params:     []constructs.Named{keyParam, valueParam},
			Package:    pkg,
		})

		// map[TKey, TValue any] interface
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			TypeParams: tp,
			Methods:    b.toNamedList(methods),
			Package:    pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `Map`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
func (b *bakerImp) BakePointer() constructs.Object {
	return bakeOnce(b, `pointer[T any]`, func() constructs.Object {
		pkg := b.BakeBuiltin()
		tp := b.proj.NewNamed(constructs.NamedArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		methods := map[string]constructs.TypeDesc{}

		// $deref() T
		methods[`$deref`] = b.proj.NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     tp,
			Package:    pkg,
		})

		// pointer[T any] interface
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    b.toNamedList(methods),
			Package:    pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `Pointer`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakeComplex64 bakes in an interface to represent a Go 64-bit complex number.
//
//	type complex64 interface {
//		$real() float32
//		$imag() float32
//	}
func (b *bakerImp) BakeComplex64() constructs.Object {
	return bakeOnce(b, `complex64`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		// func() float32
		getF := b.proj.NewSignature(constructs.SignatureArgs{
			Return:  b.BakeBasic(`float32`),
			Package: pkg,
		})

		methods := map[string]constructs.TypeDesc{
			`$real`: getF, // $real() float32
			`$imag`: getF, // $imag() float32
		}

		// complex64
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			Methods: b.toNamedList(methods),
			Package: pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `complex64`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakeComplex128 bakes in an interface to represent a Go 64-bit complex number.
//
//	type complex128 interface {
//		$real() float64
//		$imag() float64
//	}
func (b *bakerImp) BakeComplex128() constructs.Object {
	return bakeOnce(b, `complex128`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		// func() float64
		getF := b.proj.NewSignature(constructs.SignatureArgs{
			Return:  b.BakeBasic(`float64`),
			Package: pkg,
		})

		methods := map[string]constructs.TypeDesc{
			`$real`: getF, // $real() float64
			`$imag`: getF, // $imag() float64
		}

		// complex128
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			Methods: b.toNamedList(methods),
			Package: pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `complex128`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakeError bakes in an interface to represent a Go error.
//
//	type error interface {
//		Error() string
//	}
func (b *bakerImp) BakeError() constructs.Object {
	return bakeOnce(b, `error`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		// func() string
		getStr := b.proj.NewSignature(constructs.SignatureArgs{
			Return:  b.BakeBasic(`string`),
			Package: pkg,
		})

		methods := map[string]constructs.TypeDesc{
			`Error`: getStr, // Error() string
		}

		// interface { Error() string }
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			Methods: b.toNamedList(methods),
			Package: pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `error`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}

// BakeComparable bakes in an interface to represent a Go comparable.
//
//	type comparable interface {
//		$compare(other any) int
//	}
func (b *bakerImp) BakeComparable() constructs.Object {
	return bakeOnce(b, `comparable`, func() constructs.Object {
		pkg := b.BakeBuiltin()

		// other any
		otherParam := b.proj.NewNamed(constructs.NamedArgs{
			Name: `other`,
			Type: b.BakeAny(),
		})

		// func(other any) int
		getStr := b.proj.NewSignature(constructs.SignatureArgs{
			Params:  []constructs.Named{otherParam},
			Return:  b.BakeBasic(`int`),
			Package: pkg,
		})

		methods := map[string]constructs.TypeDesc{
			`$compare`: getStr, // $compare(other any) int
		}

		// interface { $compare(other any) int }
		in := b.proj.NewInterface(constructs.InterfaceArgs{
			Methods: b.toNamedList(methods),
			Package: pkg,
		})

		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `comparable`,
			Type:     in,
			Location: locs.NoLoc(),
		})
	})
}
