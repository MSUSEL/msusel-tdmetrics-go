// The baked in types are stored for quick lookup and
// to ensure only one instance of each type is created.
//
// The function names are prepended with a `$` to avoid duck-typing
// with user-defined types. Some of these types don't normally exist in Go,
// but are used to represent the construct in a way that can be abstracted.
// Other types represent the built-in types such as error.
package baker

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

const BuiltinName = `$builtin`

type Baker interface {
	ByName(name string) constructs.Construct
	BakeBuiltin() constructs.Package
	BakeAny() constructs.InterfaceDecl
	BakeList() constructs.InterfaceDecl
	BakeChan() constructs.InterfaceDecl
	BakeMap() constructs.InterfaceDecl
	BakePointer() constructs.InterfaceDecl
	BakeComplex64() constructs.InterfaceDecl
	BakeComplex128() constructs.InterfaceDecl
	BakeError() constructs.InterfaceDecl
	BakeComparable() constructs.InterfaceDecl
	BakeAppend() constructs.Method
	BakeCap() constructs.Method
	BakeClear() constructs.Method
	BakeClose() constructs.Method
	BakeComplex() constructs.Method
	BakeCopy() constructs.Method
	BakeDelete() constructs.Method
	BakeImag() constructs.Method
	BakeLen() constructs.Method
	BakeMake() constructs.Method
	BakeMax() constructs.Method
	BakeMin() constructs.Method
	BakeNew() constructs.Method
	BakePanic() constructs.Method
	BakePrint() constructs.Method
	BakePrintln() constructs.Method
	BakeReal() constructs.Method
	BakeRecover() constructs.Method
}

type bakerImp struct {
	proj  constructs.Project
	baked map[string]any
}

func New(proj constructs.Project) Baker {
	return &bakerImp{
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

// ByName gets the type by name or returns nil.
func (b *bakerImp) ByName(name string) constructs.Construct {
	switch name {
	case `error`:
		return b.BakeError()
	case `comparable`:
		return b.BakeComparable()
	case `append`:
		return b.BakeAppend()
	case `cap`:
		return b.BakeCap()
	case `clear`:
		return b.BakeClear()
	case `close`:
		return b.BakeClose()
	case `complex`:
		return b.BakeComplex()
	case `copy`:
		return b.BakeCopy()
	case `delete`:
		return b.BakeDelete()
	case `imag`:
		return b.BakeImag()
	case `len`:
		return b.BakeLen()
	case `make`:
		return b.BakeMake()
	case `max`:
		return b.BakeMax()
	case `min`:
		return b.BakeMin()
	case `new`:
		return b.BakeNew()
	case `panic`:
		return b.BakePanic()
	case `print`:
		return b.BakePrint()
	case `println`:
		return b.BakePrintln()
	case `real`:
		return b.BakeReal()
	case `recover`:
		return b.BakeRecover()
	default:
		return nil
	}
}

// BakeBuiltin bakes in a package to represent the builtin package.
func (b *bakerImp) BakeBuiltin() constructs.Package {
	return bakeOnce(b, BuiltinName, func() constructs.Package {
		builtinPkg := &packages.Package{
			PkgPath: BuiltinName,
			Name:    BuiltinName,
			Fset:    b.proj.Locs().FileSet(),
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
		pkg := b.BakeBuiltin()
		real := types.NewInterfaceType(nil, nil)

		// any interface{}
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			RealType: real,
			Package:  pkg,
			Name:     `any`,
			Exported: true,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				RealType: real,
				Package:  pkg.Source(),
			}),
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
	return bakeOnce(b, `List[T any]`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// T any
		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		tps := []constructs.TypeParam{tp}

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
		lenFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$len`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{intArg},
				Package: pkg.Source(),
			}),
		})

		// $get(index int) T
		getFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$get`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{indexArg},
				Results: []constructs.Argument{valueArg},
				Package: pkg.Source(),
			}),
		})

		// $set(index int, value T)
		setFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$set`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{indexArg, valueArg},
				Package: pkg.Source(),
			}),
		})

		// List[T]
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    pkg,
			Name:       `List`,
			Exported:   true,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{lenFunc, getFunc, setFunc},
				Package:   pkg.Source(),
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
// Note: Doesn't have `cap`, `trySend`, or `tryRecv` as defined in reflect
// because those aren't important for abstraction
func (b *bakerImp) BakeChan() constructs.InterfaceDecl {
	return bakeOnce(b, `Chan[T any]`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// T any
		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		tps := []constructs.TypeParam{tp}

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// value T
		valueArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: tp,
		})

		// okay bool
		okayArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `okay`,
			Type: b.bakeBasic(types.Bool),
		})

		// $len() int
		lenFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$len`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{intArg},
				Package: pkg.Source(),
			}),
		})

		// $recv() (T, bool)
		recvFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$recv`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{valueArg, okayArg},
				Package: pkg.Source(),
			}),
		})

		// $send(value T)
		sendFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$send`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{valueArg},
				Package: pkg.Source(),
			}),
		})

		// Chan[T]
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    pkg,
			Name:       `Chan`,
			Exported:   true,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{lenFunc, recvFunc, sendFunc},
				Package:   pkg.Source(),
			}),
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
func (b *bakerImp) BakeMap() constructs.InterfaceDecl {
	return bakeOnce(b, `Map[TKey comparable, TValue any]`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

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

		tps := []constructs.TypeParam{tpKey, tpValue}

		// <unnamed> int
		intArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// key TKey
		keyArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `key`,
			Type: tpKey,
		})

		// value TValue
		valueArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `value`,
			Type: tpValue,
		})

		// found bool
		foundArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `found`,
			Type: b.bakeBasic(types.Bool),
		})

		// $len() int
		lenFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$len`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{intArg},
				Package: pkg.Source(),
			}),
		})

		// $get(key TKey) (TValue, bool)
		getFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$get`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{keyArg},
				Results: []constructs.Argument{valueArg, foundArg},
				Package: pkg.Source(),
			}),
		})

		// $set(key TKey, value TValue)
		setFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$set`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{keyArg, valueArg},
				Package: pkg.Source(),
			}),
		})

		// Map[TKey, TValue]
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    pkg,
			Name:       `Map`,
			Exported:   true,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{lenFunc, getFunc, setFunc},
				Package:   pkg.Source(),
			}),
		})
	})
}

// BakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
func (b *bakerImp) BakePointer() constructs.InterfaceDecl {
	return bakeOnce(b, `Pointer[T any]`, func() constructs.InterfaceDecl {
		pkg := b.BakeBuiltin()

		// T any
		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})
		tps := []constructs.TypeParam{tp}

		// <unnamed> T
		resultArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: tp,
		})

		// $deref() T
		derefFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$deref`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{resultArg},
				Package: pkg.Source(),
			}),
		})

		// Pointer[T]
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:    pkg,
			Name:       `Pointer`,
			Exported:   true,
			Location:   locs.NoLoc(),
			TypeParams: tps,
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{derefFunc},
				Package:   pkg.Source(),
			}),
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
			Type: b.bakeBasic(types.Float32),
		})

		// $real() float32
		realFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$real`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{floatArg},
				Package: pkg.Source(),
			}),
		})

		// $imag() float32
		imagFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$imag`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{floatArg},
				Package: pkg.Source(),
			}),
		})

		// complex64
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `complex64`,
			Exported: true,
			Location: locs.NoLoc(),
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{realFunc, imagFunc},
				Package:   pkg.Source(),
			}),
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
			Type: b.bakeBasic(types.Float64),
		})

		// $real() float64
		realFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$real`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{floatArg},
				Package: pkg.Source(),
			}),
		})

		// $imag() float64
		imagFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$imag`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{floatArg},
				Package: pkg.Source(),
			}),
		})

		// complex128
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `complex128`,
			Exported: true,
			Location: locs.NoLoc(),
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{realFunc, imagFunc},
				Package:   pkg.Source(),
			}),
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
			Type: b.bakeBasic(types.String),
		})

		// func Error() string
		errFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `Error`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{stringArg},
				Package: pkg.Source(),
			}),
		})

		// error
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `error`,
			Exported: true,
			Location: locs.NoLoc(),
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{errFunc},
				Package:   pkg.Source(),
			}),
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
			Type: b.bakeBasic(types.Int),
		})

		// other any
		otherArg := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `other`,
			Type: b.BakeAny(),
		})

		// func $compare(other any) int
		cmpFunc := b.proj.NewAbstract(constructs.AbstractArgs{
			Name:     `$compare`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{otherArg},
				Results: []constructs.Argument{intArg},
				Package: pkg.Source(),
			}),
		})

		// comparable
		return b.proj.NewInterfaceDecl(constructs.InterfaceDeclArgs{
			Package:  pkg,
			Name:     `comparable`,
			Exported: true,
			Location: locs.NoLoc(),
			Interface: b.proj.NewInterfaceDesc(constructs.InterfaceDescArgs{
				Abstracts: []constructs.Abstract{cmpFunc},
				Package:   pkg.Source(),
			}),
		})
	})
}

// BakeAppend creates the builtin append function.
//
//	func append(slice []Type, elems ...Type) []Type
func (b *bakerImp) BakeAppend() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeCap creates the builtin cap function.
//
//	func cap(v Type) int
func (b *bakerImp) BakeCap() constructs.Method {
	return bakeOnce(b, `cap`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `v`,
			Type: b.BakeAny(),
		})
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// func cap(v any) int
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `cap`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{param},
				Results: []constructs.Argument{ret},
				Package: pkg.Source(),
			}),
		})
	})
}

// BakeClear creates the builtin clear function.
//
//	func clear[T ~[]Type | ~map[Type]Type1](t T)
func (b *bakerImp) BakeClear() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeClose creates the builtin close function.
//
//	func close(c chan<- Type)
func (b *bakerImp) BakeClose() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeComplex creates the builtin complex function.
//
//	func complex(r, i FloatType) ComplexType
func (b *bakerImp) BakeComplex() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeCopy creates the builtin copy function.
//
//	func copy(dst, src []Type) int
func (b *bakerImp) BakeCopy() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeDelete creates the builtin delete function.
//
//	func delete(m map[Type]Type1, key Type)
func (b *bakerImp) BakeDelete() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeImag creates the builtin imag function.
//
//	func imag(c ComplexType) FloatType
func (b *bakerImp) BakeImag() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeLen creates the builtin len function.
//
//	func len(v Type) int
func (b *bakerImp) BakeLen() constructs.Method {
	return bakeOnce(b, `len`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `v`,
			Type: b.BakeAny(),
		})
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// func len(v any) int
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `len`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{param},
				Results: []constructs.Argument{ret},
				Package: pkg.Source(),
			}),
		})
	})
}

// BakeMake creates the builtin make function.
//
//	func make(t Type, size ...IntegerType) Type
func (b *bakerImp) BakeMake() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeMax creates the builtin max function.
//
//	func max[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMax() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeMin creates the builtin min function.
//
//	func min[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMin() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeNew creates the builtin new function.
//
//	func new(Type) *Type
func (b *bakerImp) BakeNew() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakePanic creates the builtin panic function.
//
//	func panic(v any)
func (b *bakerImp) BakePanic() constructs.Method {
	return bakeOnce(b, `panic`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `v`,
			Type: b.BakeAny(),
		})

		// func panic(v any)
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `panic`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:  []constructs.Argument{param},
				Package: pkg.Source(),
			}),
		})
	})
}

// BakePrint creates the builtin print function.
//
//	func print(args ...Type)
func (b *bakerImp) BakePrint() constructs.Method {
	return bakeOnce(b, `print`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `args`,
			Type: b.BakeAny(),
		})

		// func print(args ...any)
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `print`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:   []constructs.Argument{ret},
				Variadic: true,
				Package:  pkg.Source(),
			}),
		})
	})
}

// BakePrintln creates the builtin println function.
//
//	func println(args ...Type)
func (b *bakerImp) BakePrintln() constructs.Method {
	return bakeOnce(b, `println`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `args`,
			Type: b.BakeAny(),
		})

		// func println(args ...any)
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `println`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Params:   []constructs.Argument{ret},
				Variadic: true,
				Package:  pkg.Source(),
			}),
		})
	})
}

// BakeReal creates the builtin real function.
//
//	func real(c ComplexType) FloatType
func (b *bakerImp) BakeReal() constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeRecover creates the builtin recover function.
//
//	func recover() any
func (b *bakerImp) BakeRecover() constructs.Method {
	return bakeOnce(b, `recover`, func() constructs.Method {
		pkg := b.BakeBuiltin()
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeAny(),
		})

		// func recover() any
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:  pkg,
			Name:     `recover`,
			Exported: true,
			Signature: b.proj.NewSignature(constructs.SignatureArgs{
				Results: []constructs.Argument{ret},
				Package: pkg.Source(),
			}),
		})
	})
}
