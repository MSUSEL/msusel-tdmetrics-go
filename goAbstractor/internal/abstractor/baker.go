// The baked in types are stored for quick lookup and
// to ensure only one instance of each type is created.
//
// The function names are prepended with a `$` to avoid duck-typing
// with user-defined types. These types don't normally exist in Go,
// but are used to represent the construct in a way that can be abstracted.
package abstractor

import (
	"fmt"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func bakeOnce[T any](ab *abstractor, key string, create func() T) T {
	if baked, has := ab.baked[key]; has {
		t, ok := baked.(T)
		if !ok {
			panic(fmt.Errorf(`unexpected type for %[1]s: wanted %[2]T, got %[3]T: %[3]v`, key, utils.Zero[T](), baked))
		}
		return t
	}

	t := create()
	ab.baked[key] = t
	return t
}

// bakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (ab *abstractor) bakeAny() constructs.Interface {
	return bakeOnce(ab, `any`, func() constructs.Interface {
		// any
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			RealType: types.NewInterfaceType(nil, nil),
		})
	})
}

// bakeIntFunc bakes in a signature for `func() int`.
// This is useful for things like `cap() int` and `len() int`.
func (ab *abstractor) bakeIntFunc() constructs.Signature {
	return bakeOnce(ab, `func() int`, func() constructs.Signature {
		// func() int
		return constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			Return: constructs.BasicFor[int](ab.proj.Types()),
		})
	})
}

// bakeReturnTuple bakes in a structure used for a return value
// tuple with a variable type value and a boolean.
//
//	struct {
//		value T
//		ok    bool
//	}
func (ab *abstractor) bakeReturnTuple(tp constructs.Named) constructs.Struct {
	return bakeOnce(ab, `struct[T] { value T; ok bool }`, func() constructs.Struct {
		fieldValue := constructs.NewNamed(ab.proj.Types(), `value`, tp)
		fieldOk := constructs.NewNamed(ab.proj.Types(), `ok`, constructs.BasicFor[bool](ab.proj.Types()))

		// struct
		return constructs.NewStruct(ab.proj.Types(), constructs.StructArgs{
			Fields: []constructs.Named{fieldValue, fieldOk},
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
func (ab *abstractor) bakeList() constructs.Interface {
	return bakeOnce(ab, `list[T any]`, func() constructs.Interface {
		tp := constructs.NewNamed(ab.proj.Types(), `T`, ab.bakeAny())

		intFunc := ab.bakeIntFunc()
		indexParam := constructs.NewNamed(ab.proj.Types(), `index`, constructs.BasicFor[int](ab.proj.Types()))
		valueParam := constructs.NewNamed(ab.proj.Types(), `value`, tp)

		methods := map[string]constructs.TypeDesc{}
		methods[`$len`] = intFunc // $len() int
		methods[`$cap`] = intFunc // $cap() int

		// $get(index int) T
		methods[`$get`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{indexParam},
			Return:     tp,
		})

		// $set(index int, value T)
		methods[`$set`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{indexParam, valueParam},
		})

		// list[T any] interface
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			TypeParams:   []constructs.Named{tp},
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeChan bakes in an interface to represent a Go chan:
//
//	type chan[T any] interface {
//		$len() int
//		$recv() (T, bool)
//		$send(value T)
//	}
//
// Note: Doesn't currently have cap, trySend, or tryRecv as defined in reflect.
func (ab *abstractor) bakeChan() constructs.Interface {
	return bakeOnce(ab, `chan[T any]`, func() constructs.Interface {
		tp := constructs.NewNamed(ab.proj.Types(), `T`, ab.bakeAny())
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = ab.bakeIntFunc()

		// $recv() (T, bool)
		methods[`$recv`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     ab.bakeReturnTuple(tp),
		})

		// $send(value T)
		methods[`$send`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{constructs.NewNamed(ab.proj.Types(), `value`, tp)},
		})

		// chan[T any] interface
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			TypeParams:   []constructs.Named{tp},
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeMap bakes in an interface to represent a Go map:
//
//	type map[TKey, TValue any] interface {
//		$len() int
//		$get(key TKey) (TValue, bool)
//		$set(key TKey, value TValue)
//	}
//
// Note: Doesn't currently require Key to be comparable as defined in reflect.
func (ab *abstractor) bakeMap() constructs.Interface {
	return bakeOnce(ab, `map[TKey, TValue any]`, func() constructs.Interface {
		tpKey := constructs.NewNamed(ab.proj.Types(), `TKey`, ab.bakeAny())
		tpValue := constructs.NewNamed(ab.proj.Types(), `TValue`, ab.bakeAny())
		tp := []constructs.Named{tpKey, tpValue}
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = ab.bakeIntFunc()

		// $get(key TKey) (TValue, bool)
		methods[`$get`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: tp,
			Params:     []constructs.Named{constructs.NewNamed(ab.proj.Types(), `key`, tpKey)},
			Return:     ab.bakeReturnTuple(tpValue),
		})

		// $set(key TKey, value TValue)
		methods[`$set`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: tp,
			Params: []constructs.Named{
				constructs.NewNamed(ab.proj.Types(), `key`, tpKey),
				constructs.NewNamed(ab.proj.Types(), `value`, tpValue),
			},
		})

		// map[TKey, TValue any] interface
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			TypeParams:   tp,
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
func (ab *abstractor) bakePointer() constructs.Interface {
	return bakeOnce(ab, `pointer[T any]`, func() constructs.Interface {
		tp := constructs.NewNamed(ab.proj.Types(), `T`, ab.bakeAny())
		methods := map[string]constructs.TypeDesc{}

		// $deref() T
		methods[`$deref`] = constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     tp,
		})

		// pointer[T any] interface
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			TypeParams:   []constructs.Named{tp},
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeComplex64 bakes in an interface to represent a Go 64-bit complex number.
func (ab *abstractor) bakeComplex64() constructs.Interface {
	return bakeOnce(ab, `complex64`, func() constructs.Interface {

		// func() float32
		getF := constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			Return: constructs.BasicFor[float32](ab.proj.Types()),
		})

		methods := map[string]constructs.TypeDesc{}
		methods[`$real`] = getF // $real() float32
		methods[`$imag`] = getF // $imag() float32

		// complex64
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeComplex128 bakes in an interface to represent a Go 64-bit complex number.
func (ab *abstractor) bakeComplex128() constructs.Interface {
	return bakeOnce(ab, `complex128`, func() constructs.Interface {

		// func() float64
		getF := constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			Return: constructs.BasicFor[float64](ab.proj.Types()),
		})

		methods := map[string]constructs.TypeDesc{}
		methods[`$real`] = getF // $real() float64
		methods[`$imag`] = getF // $imag() float64

		// complex128
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeError bakes in an interface to represent a Go error.
func (ab *abstractor) bakeError() constructs.Interface {
	return bakeOnce(ab, `error`, func() constructs.Interface {

		// func() string
		getStr := constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			Return: constructs.BasicFor[string](ab.proj.Types()),
		})

		methods := map[string]constructs.TypeDesc{}
		methods[`Error`] = getStr // Error() string

		// interface { Error() string }
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeComparable bakes in an interface to represent a Go comparable.
func (ab *abstractor) bakeComparable() constructs.Interface {
	return bakeOnce(ab, `comparable`, func() constructs.Interface {
		tp := constructs.NewNamed(ab.proj.Types(), `T`, ab.bakeAny())

		// func(other T) int
		getStr := constructs.NewSignature(ab.proj.Types(), constructs.SignatureArgs{
			Params: []constructs.Named{
				constructs.NewNamed(ab.proj.Types(), `other`, tp),
			},
			Return: constructs.BasicFor[int](ab.proj.Types()),
		})

		methods := map[string]constructs.TypeDesc{}
		methods[`$compare`] = getStr // $compare(other T) int

		// interface { $compare(other T) int }
		return constructs.NewInterface(ab.proj.Types(), constructs.InterfaceArgs{
			Methods:      methods,
			InitInherits: []constructs.Interface{ab.bakeAny()},
		})
	})
}

// bakeBuiltin bakes in a package to represent the builtin package.
func (ab *abstractor) bakeBuiltin() constructs.Package {
	return bakeOnce(ab, `$builtin`, func() constructs.Package {
		errTyp := constructs.NewTypeDef(`error`, ab.bakeError())
		cmpTyp := constructs.NewTypeDef(`comparable`, ab.bakeComparable())

		pkg := constructs.NewPackage(nil, `$builtin`, `$builtin`, []string{})
		pkg.AppendTypes(errTyp, cmpTyp)
		return pkg
	})
}
