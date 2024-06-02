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

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

// TODO: finish updating to use bakeOnce
func bakeOnce[T typeDesc.TypeDesc](ab *abstractor, key string, create func() T) T {
	if baked, has := ab.baked[key]; has {
		t, ok := baked.(T)
		if ok {
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
func (ab *abstractor) bakeAny() typeDesc.Interface {
	return bakeOnce(ab, `any`, func() typeDesc.Interface {
		// any
		return typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
			RealType: types.NewInterfaceType(nil, nil),
		})
	})
}

// bakeIntFunc bakes in a signature for `func() int`.
// This is useful for things like `cap() int` and `len() int`.
func (ab *abstractor) bakeIntFunc() typeDesc.Signature {
	return bakeOnce(ab, `func() int`, func() typeDesc.Signature {
		// func() int
		return typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
			Return: typeDesc.BasicFor[int](ab.proj),
		})
	})
}

// bakeReturnTuple bakes in a structure used for a return value
// tuple with a variable type value and a boolean.
//
//	struct[T any] {
//		value T
//		ok    bool
//	}
func (ab *abstractor) bakeReturnTuple() typeDesc.Struct {
	return bakeOnce(ab, `struct[T] { value T; ok bool }`, func() typeDesc.Struct {
		tp := typeDesc.NewNamed(ab.proj, `T`, ab.bakeAny())
		fieldValue := typeDesc.NewNamed(ab.proj, `value`, tp)
		fieldOk := typeDesc.NewNamed(ab.proj, `ok`, typeDesc.BasicFor[bool](ab.proj))
		// struct[T any]s
		return typeDesc.NewStruct(ab.proj, typeDesc.StructArgs{
			TypeParams: []typeDesc.Named{tp},
			Fields:     []typeDesc.Named{fieldValue, fieldOk},
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
func (ab *abstractor) bakeList() typeDesc.Interface {
	const bakeKey = `list[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	tp := typeDesc.NewNamed(ab.proj, `T`, ab.bakeAny())

	intFunc := ab.bakeIntFunc()
	indexParam := typeDesc.NewNamed(ab.proj, `index`, typeDesc.BasicFor[int](ab.proj))
	valueParam := typeDesc.NewNamed(ab.proj, `value`, tp)

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$len`] = intFunc // $len() int
	methods[`$cap`] = intFunc // $cap() int

	// $get(index int) T
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{indexParam},
		Return:     tp,
	})
	methods[`$get`] = typeDesc.NewSolid(ab.proj, nil, getF, tp)

	// $set(index int, value T)
	setF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{indexParam, valueParam},
	})
	methods[`$set`] = typeDesc.NewSolid(ab.proj, nil, setF, tp)

	// list[T any] interface
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
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
func (ab *abstractor) bakeChan() typeDesc.Interface {
	const bakeKey = `chan[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	tp := typeDesc.NewNamed(ab.proj, `T`, ab.bakeAny())
	methods := map[string]typeDesc.TypeDesc{}

	// $len() int
	methods[`$len`] = ab.bakeIntFunc()

	// $recv() (T, bool)
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Return:     typeDesc.NewSolid(ab.proj, nil, ab.bakeReturnTuple(), tp),
	})
	methods[`$recv`] = typeDesc.NewSolid(ab.proj, nil, getF, tp)

	// $send(value T)
	setF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{typeDesc.NewNamed(ab.proj, `value`, tp)},
	})
	methods[`$send`] = typeDesc.NewSolid(ab.proj, nil, setF, tp)

	// chan[T any] interface
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
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
func (ab *abstractor) bakeMap() typeDesc.Interface {
	const bakeKey = `map[TKey, TValue any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	tpKey := typeDesc.NewNamed(ab.proj, `TKey`, ab.bakeAny())
	tpValue := typeDesc.NewNamed(ab.proj, `TValue`, ab.bakeAny())
	tp := []typeDesc.Named{tpKey, tpValue}
	methods := map[string]typeDesc.TypeDesc{}

	// $len() int
	methods[`$len`] = ab.bakeIntFunc()

	// $get(key TKey) (TValue, bool)
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: tp,
		Params:     []typeDesc.Named{typeDesc.NewNamed(ab.proj, `key`, tpKey)},
		Return:     typeDesc.NewSolid(ab.proj, nil, ab.bakeReturnTuple(), tpValue),
	})
	methods[`$get`] = typeDesc.NewSolid(ab.proj, nil, getF, tpKey, tpValue)

	// $set(key TKey, value TValue)
	setF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: tp,
		Params: []typeDesc.Named{
			typeDesc.NewNamed(ab.proj, `key`, tpKey),
			typeDesc.NewNamed(ab.proj, `value`, tpValue),
		},
	})
	methods[`$set`] = typeDesc.NewSolid(ab.proj, nil, setF, tpKey, tpValue)

	// map[TKey, TValue any] interface
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		TypeParams:   tp,
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
}

// bakePointer bakes in an interface to represent a Go pointer:
//
//	type pointer[T any] interface {
//		$deref() T
//	}
func (ab *abstractor) bakePointer() typeDesc.Interface {
	const bakeKey = `pointer[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	tp := typeDesc.NewNamed(ab.proj, `T`, ab.bakeAny())
	methods := map[string]typeDesc.TypeDesc{}

	// $deref() T
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Return:     tp,
	})
	methods[`$deref`] = typeDesc.NewSolid(ab.proj, nil, getF, tp)

	// pointer[T any] interface
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
}

// bakeComplex64 bakes in an interface to represent a Go 64-bit complex number.
func (ab *abstractor) bakeComplex64() typeDesc.Interface {
	const bakeKey = `complex64`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	// func() float32
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		Return: typeDesc.BasicFor[float32](ab.proj),
	})

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$real`] = getF // $real() float32
	methods[`$imag`] = getF // $imag() float32

	// complex64
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
}

// bakeComplex128 bakes in an interface to represent a Go 64-bit complex number.
func (ab *abstractor) bakeComplex128() typeDesc.Interface {
	const bakeKey = `complex128`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	// func() float64
	getF := typeDesc.NewSignature(ab.proj, typeDesc.SignatureArgs{
		Return: typeDesc.BasicFor[float64](ab.proj),
	})

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$real`] = getF // $real() float32
	methods[`$imag`] = getF // $imag() float32

	// complex128
	t := typeDesc.NewInterface(ab.proj, typeDesc.InterfaceArgs{
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	ab.baked[bakeKey] = t
	return t
}
