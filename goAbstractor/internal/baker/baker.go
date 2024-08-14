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

type Baker interface {
	BakeBuiltin() constructs.Package
	BakeAny() constructs.InterfaceDecl
	BakeList(elem constructs.TypeDesc) constructs.InterfaceDecl
	BakeChan(elem constructs.TypeDesc) constructs.InterfaceDecl
	BakeMap(key, value constructs.TypeDesc) constructs.InterfaceDecl
	BakePointer(elem constructs.TypeDesc) constructs.InterfaceDecl
	BakeComplex64() constructs.InterfaceDecl
	BakeComplex128() constructs.InterfaceDecl
	BakeError() constructs.InterfaceDecl
	BakeComparable() constructs.InterfaceDecl
}

type bakerImp struct {
	fSet  *token.FileSet
	proj  constructs.Project
	baked map[string]any
}

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

// BakeBuiltin bakes in a package to represent the builtin package.
func (b *bakerImp) BakeBuiltin() constructs.Package {
	return bakeOnce(b, BuiltinName, func() constructs.Package {
		builtinPkg := &packages.Package{
			PkgPath: BuiltinName,
			Name:    BuiltinName,
			Fset:    b.fSet,
			Types:   types.NewPackage(BuiltinName, BuiltinName),
		}

		// package $builtin
		return b.proj.NewPackage(constructs.PackageArgs{
			RealPkg: builtinPkg,
			Path:    BuiltinName,
			Name:    BuiltinName,
		})
	})
}

// bakeBasic bakes in a basic type by name. The name must
// be a valid basic type (e.g. int, int32, float64)
// but may not be complex numbers or interfaces like any.
func (b *bakerImp) bakeBasic(kind types.BasicKind) constructs.Basic {
	bk := types.Typ[kind]
	return bakeOnce(b, `basic `+bk.Name(), func() constructs.Basic {
		return b.proj.NewBasic(constructs.BasicArgs{
			RealType: bk,
		})
	})
}

// BakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (b *bakerImp) BakeAny() constructs.InterfaceDecl {
	return bakeOnce(b, `any`, func() constructs.InterfaceDecl {
		// any
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			RealType: types.NewInterfaceType(nil, nil),
			Package:  b.BakeBuiltin(),
			Name:     `any`,
			Location: locs.NoLoc(),
		})
	})
}

// bakeList bakes in an interface to represent a Go array or slice:
//
//	type list[T any] interface {
//		$len() int
//		$get(index int) T
//		$set(index int, value T)
//	}
//
// Note: The difference between an array and slice aren't
// important for abstraction, so they are combined into one.
// Also `cap` and `offset` aren't important, so ignored.
func (b *bakerImp) BakeList() constructs.InterfaceDecl {
	return bakeOnce(b, `List`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()
		var tps []constructs.TypeParam

		// T any
		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		tps = []constructs.TypeParam{tp}

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// index int
		indexArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `index`,
			Type: b.bakeBasic(types.Int),
		})

		// value T
		valueArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: tp,
		})

		// $len() int
		lenFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$len`,
			Location: locs.NoLoc(),
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{intArg},
			}),
		})

		// $get(index int) T
		getFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$get`,
			Location: locs.NoLoc(),
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{indexArg},
				Results: []constructs.Argument{valueArg},
			}),
		})

		// $set(index int, value T)
		setFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$get`,
			Location: locs.NoLoc(),
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params: []constructs.Argument{indexArg, valueArg},
			}),
		})

		// List[T]
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    pkg,
			Name:       `List`,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{

				Methods: []constructs.Method{lenFunc, getFunc, setFunc},
			}),
		})
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
// Note: Doesn't have `cap`, `trySend`, or `tryRecv` as defined in reflect
// because those aren't important for abstraction
func (b *bakerImp) BakeChan(elem constructs.TypeDesc) constructs.InterfaceDecl {
	generic := utils.IsNil(elem)
	bakeKey := `Chan[T any]`
	if !generic {
		bakeKey = `Chan@[` + elem.GoType().String() + `]`
	}
	return bakeOnce(b, bakeKey, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()
		var tps []constructs.TypeParam

		if generic {
			// T any
			tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
				Name: `T`,
				Type: b.BakeAny(),
			})
			tps = []constructs.TypeParam{tp}
			elem = tp
		}

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`int`),
		})

		// value T
		valueArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: elem,
		})

		// okay bool
		okayArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `okay`,
			Type: b.BakeBasic(`bool`),
		})

		// $len() int
		lenFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$len`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{intArg},
		})

		// $recv() (T, bool)
		recvFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$recv`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{valueArg, okayArg},
		})

		// $send(value T)
		sendFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$send`,
			Location: locs.NoLoc(),
			Params:   []constructs.Argument{valueArg},
		})

		// Chan[T]
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:    pkg,
			Name:       `Chan`,
			TypeParams: tps,
			Methods:    []constructs.Method{lenFunc, recvFunc, sendFunc},
			Location:   locs.NoLoc(),
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
// If the given key and value are nil, then the generic form is returned.
// Otherwise, the instance realization on the given key and value is returned.
//
// Note: Doesn't currently require Key to be comparable as defined in reflect.
func (b *bakerImp) BakeMap(key, value constructs.TypeDesc) constructs.InterfaceDecl {
	generic := utils.IsNil(key)
	if utils.IsNil(value) != generic {
		panic(terror.New(`instance of map must have both key and value not nil, otherwise both nil`).
			With(`key`, key).
			With(`value`, value))
	}
	bakeKey := `Map[TKey comparable, TValue any]`
	if !generic {
		bakeKey = `Chan@[` + key.GoType().String() + `, ` + value.GoType().String() + `]`
	}
	return bakeOnce(b, bakeKey, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()
		var tps []constructs.TypeParam

		if generic {
			// TKey comparable
			tpKey := b.proj.NewTypeParam(constructs.TypeParamArgs{
				Name: `TKey`,
				Type: b.BakeComparable(),
			})

			// TValue any
			tpValue := b.proj.NewTypeParam(constructs.TypeParamArgs{
				Name: `TValue`,
				Type: b.BakeAny(),
			})

			tps = []constructs.TypeParam{tpKey, tpValue}
			key, value = tpKey, tpValue
		}

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`int`),
		})

		// key TKey
		keyArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `key`,
			Type: key,
		})

		// value TValue
		valueArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: value,
		})

		// found bool
		foundArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `found`,
			Type: b.BakeBasic(`bool`),
		})

		// $len() int
		lenFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$len`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{intArg},
		})

		// $get(key TKey) (TValue, bool)
		getFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$get`,
			Location: locs.NoLoc(),
			Params:   []constructs.Argument{keyArg},
			Results:  []constructs.Argument{valueArg, foundArg},
		})

		// $set(key TKey, value TValue)
		setFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$set`,
			Location: locs.NoLoc(),
			Params:   []constructs.Argument{keyArg, valueArg},
		})

		// Map[TKey, TValue]
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:    pkg,
			Name:       `Map`,
			TypeParams: tps,
			Methods:    []constructs.Method{lenFunc, getFunc, setFunc},
			Location:   locs.NoLoc(),
		})
	})
}

// BakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
//
// If the given elements is nil, then the generic form is returned.
// Otherwise, the instance realization on the given element is returned.
func (b *bakerImp) BakePointer(elem constructs.TypeDesc) constructs.InterfaceDecl {
	generic := utils.IsNil(elem)
	bakeKey := `Pointer[T any]`
	if !generic {
		bakeKey = `Pointer@[` + elem.GoType().String() + `]`
	}
	return bakeOnce(b, bakeKey, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()
		var tps []constructs.TypeParam

		if generic {
			// T any
			tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
				Name: `T`,
				Type: b.BakeAny(),
			})
			tps = []constructs.TypeParam{tp}
			elem = tp
		}

		// <unnamed> T
		resultArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: elem,
		})

		// $deref() T
		derefFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       `$deref`,
			TypeParams: tps,
			Results:    []constructs.Argument{resultArg},
			Location:   locs.NoLoc(),
		})

		// Pointer[T]
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:    pkg,
			Name:       `Pointer`,
			TypeParams: tps,
			Methods:    []constructs.Method{derefFunc},
			Location:   locs.NoLoc(),
		})
	})
}

// BakeComplex64 bakes in an interface to represent a Go 64-bit complex number.
//
//	type complex64 interface {
//		$real() float32
//		$imag() float32
//	}
func (b *bakerImp) BakeComplex64() constructs.InterfaceDecl {
	return bakeOnce(b, `complex64`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// <unnamed> float32
		floatArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`float32`),
		})

		// $real() float32
		realFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$real`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{floatArg},
		})

		// $imag() float32
		imagFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$imag`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{floatArg},
		})

		// complex64
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:  pkg,
			Name:     `complex64`,
			Methods:  []constructs.Method{realFunc, imagFunc},
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
func (b *bakerImp) BakeComplex128() constructs.InterfaceDecl {
	return bakeOnce(b, `complex128`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// <unnamed> float64
		floatArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`float64`),
		})

		// $real() float64
		realFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$real`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{floatArg},
		})

		// $imag() float64
		imagFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$imag`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{floatArg},
		})

		// complex128
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:  pkg,
			Name:     `complex128`,
			Methods:  []constructs.Method{realFunc, imagFunc},
			Location: locs.NoLoc(),
		})
	})
}

// BakeError bakes in an interface to represent a Go error.
//
//	type error interface {
//		Error() string
//	}
func (b *bakerImp) BakeError() constructs.InterfaceDecl {
	return bakeOnce(b, `error`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// <unnamed> string
		stringArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeBasic(`string`),
		})

		// func Error() string
		errFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `Error`,
			Location: locs.NoLoc(),
			Results:  []constructs.Argument{stringArg},
		})

		// error
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:  pkg,
			Name:     `error`,
			Methods:  []constructs.Method{errFunc},
			Location: locs.NoLoc(),
		})
	})
}

// BakeComparable bakes in an interface to represent a Go comparable.
//
//	type comparable interface {
//		$compare(other any) int
//	}
func (b *bakerImp) BakeComparable() constructs.InterfaceDecl {
	return bakeOnce(b, `comparable`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(`int`),
		})

		// other any
		otherArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `other`,
			Type: b.BakeAny(),
		})

		// func $compare(other any) int
		cmpFunc := b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `$compare`,
			Location: locs.NoLoc(),
			Params:   []constructs.Argument{otherArg},
			Results:  []constructs.Argument{intArg},
		})

		// comparable
		return b.proj.NewInterface(constructs.InterfaceArgs{
			Package:  pkg,
			Name:     `comparable`,
			Methods:  []constructs.Method{cmpFunc},
			Location: locs.NoLoc(),
		})
	})
}
