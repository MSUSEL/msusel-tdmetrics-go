// The baked in types are stored for quick lookup and
// to ensure only one instance of each type is created.
//
// The function names are prepended with a `$` to avoid duck-typing
// with user-defined types. Some of these types don't normally exist in Go,
// but are used to represent the construct in a way that can be abstracted.
// Other types represent the built-in types such as error.
package abstractor

import (
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
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

func (ab *abstractor) toNamedList(m map[string]constructs.TypeDesc) []constructs.Named {
	n := make([]constructs.Named, len(m))
	names := utils.SortedKeys(m)
	for i, name := range names {
		n[i] = ab.proj.Types().NewNamed(name, m[name])
	}
	return n
}

// bakeBuiltin bakes in a package to represent the builtin package.
func (ab *abstractor) bakeBuiltin() constructs.Package {
	return bakeOnce(ab, `$builtin`, func() constructs.Package {
		pkg := constructs.NewPackage(constructs.PackageArgs{
			RealPkg: ab.builtin,
			Path:    `$builtin`,
			Name:    `$builtin`,
		})

		ab.proj.AppendPackage(pkg)
		return pkg
	})
}

// bakeBasic bakes in a basic type by name. The name must
// be a valid basic type (e.g. int, int32, float64)
// but may not be complex numbers or interfaces like any.
func (ab *abstractor) bakeBasic(typeName string) constructs.Named {
	return bakeOnce(ab, `basic `+typeName, func() constructs.Named {
		basic := ab.proj.Types().NewBasicFromName(ab.builtin, typeName)

		named := ab.proj.Types().NewNamed(typeName, basic)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (ab *abstractor) bakeAny() constructs.Named {
	return bakeOnce(ab, `any`, func() constructs.Named {
		// any
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			RealType: types.NewInterfaceType(nil, nil),
			Package:  ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`any`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeIntFunc bakes in a signature for `func() int`.
// This is useful for things like `cap() int` and `len() int`.
func (ab *abstractor) bakeIntFunc() constructs.Signature {
	return bakeOnce(ab, `func() int`, func() constructs.Signature {
		// func() int
		return ab.proj.Types().NewSignature(constructs.SignatureArgs{
			Return:  ab.bakeBasic(`int`),
			Package: ab.builtin,
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
	return bakeOnce(ab, `struct { value T; ok bool }`, func() constructs.Struct {
		fieldValue := ab.proj.Types().NewNamed(`value`, tp)
		fieldOk := ab.proj.Types().NewNamed(`ok`, ab.bakeBasic(`bool`))

		// struct
		return ab.proj.Types().NewStruct(constructs.StructArgs{
			Fields:  []constructs.Named{fieldValue, fieldOk},
			Package: ab.builtin,
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
func (ab *abstractor) bakeList() constructs.Named {
	return bakeOnce(ab, `list[T any]`, func() constructs.Named {
		tp := ab.proj.Types().NewNamed(`T`, ab.bakeAny())

		intFunc := ab.bakeIntFunc()
		indexParam := ab.proj.Types().NewNamed(`index`, ab.bakeBasic(`int`))
		valueParam := ab.proj.Types().NewNamed(`value`, tp)

		methods := map[string]constructs.TypeDesc{}
		methods[`$len`] = intFunc // $len() int
		methods[`$cap`] = intFunc // $cap() int

		// $get(index int) T
		methods[`$get`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{indexParam},
			Return:     tp,
			Package:    ab.builtin,
		})

		// $set(index int, value T)
		methods[`$set`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{indexParam, valueParam},
			Package:    ab.builtin,
		})

		// list[T any] interface
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    ab.toNamedList(methods),
			Package:    ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`List`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
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
func (ab *abstractor) bakeChan() constructs.Named {
	return bakeOnce(ab, `chan[T any]`, func() constructs.Named {
		tp := ab.proj.Types().NewNamed(`T`, ab.bakeAny())
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = ab.bakeIntFunc()

		// $recv() (T, bool)
		methods[`$recv`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     ab.bakeReturnTuple(tp),
			Package:    ab.builtin,
		})

		// $send(value T)
		methods[`$send`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Params:     []constructs.Named{ab.proj.Types().NewNamed(`value`, tp)},
			Package:    ab.builtin,
		})

		// chan[T any] interface
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    ab.toNamedList(methods),
			Package:    ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`Chan`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
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
func (ab *abstractor) bakeMap() constructs.Named {
	return bakeOnce(ab, `map[TKey, TValue any]`, func() constructs.Named {
		tpKey := ab.proj.Types().NewNamed(`TKey`, ab.bakeAny())
		tpValue := ab.proj.Types().NewNamed(`TValue`, ab.bakeAny())
		tp := []constructs.Named{tpKey, tpValue}
		methods := map[string]constructs.TypeDesc{}

		// $len() int
		methods[`$len`] = ab.bakeIntFunc()

		// $get(key TKey) (TValue, bool)
		methods[`$get`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: tp,
			Params:     []constructs.Named{ab.proj.Types().NewNamed(`key`, tpKey)},
			Return:     ab.bakeReturnTuple(tpValue),
			Package:    ab.builtin,
		})

		// $set(key TKey, value TValue)
		methods[`$set`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: tp,
			Params: []constructs.Named{
				ab.proj.Types().NewNamed(`key`, tpKey),
				ab.proj.Types().NewNamed(`value`, tpValue),
			},
			Package: ab.builtin,
		})

		// map[TKey, TValue any] interface
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			TypeParams: tp,
			Methods:    ab.toNamedList(methods),
			Package:    ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`Map`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
func (ab *abstractor) bakePointer() constructs.Named {
	return bakeOnce(ab, `pointer[T any]`, func() constructs.Named {
		tp := ab.proj.Types().NewNamed(`T`, ab.bakeAny())
		methods := map[string]constructs.TypeDesc{}

		// $deref() T
		methods[`$deref`] = ab.proj.Types().NewSignature(constructs.SignatureArgs{
			TypeParams: []constructs.Named{tp},
			Return:     tp,
			Package:    ab.builtin,
		})

		// pointer[T any] interface
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    ab.toNamedList(methods),
			Package:    ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`Pointer`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeComplex64 bakes in an interface to represent a Go 64-bit complex number.
//
//	type complex64 interface {
//		$real() float32
//		$imag() float32
//	}
func (ab *abstractor) bakeComplex64() constructs.Named {
	return bakeOnce(ab, `complex64`, func() constructs.Named {

		// func() float32
		getF := ab.proj.Types().NewSignature(constructs.SignatureArgs{
			Return:  ab.bakeBasic(`float32`),
			Package: ab.builtin,
		})

		methods := map[string]constructs.TypeDesc{
			`$real`: getF, // $real() float32
			`$imag`: getF, // $imag() float32
		}

		// complex64
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			Methods: ab.toNamedList(methods),
			Package: ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`complex64`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeComplex128 bakes in an interface to represent a Go 64-bit complex number.
//
//	type complex128 interface {
//		$real() float64
//		$imag() float64
//	}
func (ab *abstractor) bakeComplex128() constructs.Named {
	return bakeOnce(ab, `complex128`, func() constructs.Named {

		// func() float64
		getF := ab.proj.Types().NewSignature(constructs.SignatureArgs{
			Return:  ab.bakeBasic(`float64`),
			Package: ab.builtin,
		})

		methods := map[string]constructs.TypeDesc{
			`$real`: getF, // $real() float64
			`$imag`: getF, // $imag() float64
		}

		// complex128
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			Methods: ab.toNamedList(methods),
			Package: ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`complex128`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeError bakes in an interface to represent a Go error.
//
//	type error interface {
//		Error() string
//	}
func (ab *abstractor) bakeError() constructs.Named {
	return bakeOnce(ab, `error`, func() constructs.Named {

		// func() string
		getStr := ab.proj.Types().NewSignature(constructs.SignatureArgs{
			Return:  ab.bakeBasic(`string`),
			Package: ab.builtin,
		})

		methods := map[string]constructs.TypeDesc{
			`Error`: getStr, // Error() string
		}

		// interface { Error() string }
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			Methods: ab.toNamedList(methods),
			Package: ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`error`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}

// bakeComparable bakes in an interface to represent a Go comparable.
//
//	type comparable interface {
//		$compare(other T) int
//	}
func (ab *abstractor) bakeComparable() constructs.Named {
	return bakeOnce(ab, `comparable`, func() constructs.Named {
		tp := ab.proj.Types().NewNamed(`T`, ab.bakeAny())

		// func(other T) int
		getStr := ab.proj.Types().NewSignature(constructs.SignatureArgs{
			Params: []constructs.Named{
				ab.proj.Types().NewNamed(`other`, tp),
			},
			Return:  ab.bakeBasic(`int`),
			Package: ab.builtin,
		})

		methods := map[string]constructs.TypeDesc{}
		methods[`$compare`] = getStr // $compare(other T) int

		// interface { $compare(other T) int }
		in := ab.proj.Types().NewInterface(constructs.InterfaceArgs{
			TypeParams: []constructs.Named{tp},
			Methods:    ab.toNamedList(methods),
			Package:    ab.builtin,
		})

		named := ab.proj.Types().NewNamed(`comparable`, in)
		ab.bakeBuiltin().AppendValues(named)
		return named
	})
}
