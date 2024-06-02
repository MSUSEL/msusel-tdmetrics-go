// The baked in types are stored for quick lookup and
// to ensure only one instance of each type is created.
//
// The function names are prepended with a `$` to avoid duck-typing
// with user-defined types. These types don't normally exist in Go,
// but are used to represent the construct in a way that can be abstracted.
package abstractor

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
)

// bakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (ab *abstractor) bakeAny() typeDesc.Interface {
	const bakeKey = `any`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Interface)
	}

	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		RealType: types.NewInterfaceType(nil, nil),
	})
	t = ab.proj.RegisterInterface(t)
	ab.baked[bakeKey] = t
	return t
}

// bakeIntFunc bakes in a signature for `func() int`.
// This is useful for things like `cap() int` and `len() int`.
func (ab *abstractor) bakeIntFunc() typeDesc.Signature {
	const bakeKey = `func() int`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Signature)
	}

	// func() int
	f := typeDesc.NewSignature(typeDesc.SignatureArgs{
		Return: typeDesc.BasicFor[int](ab.proj),
	})
	f = ab.proj.RegisterSignature(f)
	ab.baked[bakeKey] = f
	return f
}

// bakeReturnTuple bakes in a structure used for a return value
// tuple with a variable type value and a boolean.
//
//	struct[T any] {
//		value T
//		ok    bool
//	}
func (ab *abstractor) bakeReturnTuple() typeDesc.Struct {
	const bakeKey = `struct[T] { value T; ok bool }`
	if t, has := ab.baked[bakeKey]; has {
		return t.(typeDesc.Struct)
	}

	tp := typeDesc.NewNamed(`T`, ab.bakeAny())
	t := typeDesc.NewStruct(typeDesc.StructArgs{
		TypeParams: []typeDesc.Named{tp},
		Fields: []typeDesc.Named{
			typeDesc.NewNamed(`value`, tp),
			typeDesc.NewNamed(`ok`, typeDesc.BasicFor[bool](ab.proj)),
		},
	})
	t = ab.proj.RegisterStruct(t)
	ab.baked[bakeKey] = t
	return t
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

	tp := typeDesc.NewNamed(`T`, ab.bakeAny())

	intFunc := ab.bakeIntFunc()
	indexParam := typeDesc.NewNamed(`index`, typeDesc.BasicFor[int](ab.proj))
	valueParam := typeDesc.NewNamed(`value`, tp)

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$len`] = intFunc // $len() int
	methods[`$cap`] = intFunc // $cap() int

	// $get(index int) T
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{indexParam},
		Return:     tp,
	})
	getF = ab.proj.RegisterSignature(getF)
	methods[`$get`] = typeDesc.NewSolid(nil, getF, tp)

	// $set(index int, value T)
	setF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{indexParam, valueParam},
	})
	setF = ab.proj.RegisterSignature(setF)
	methods[`$set`] = typeDesc.NewSolid(nil, setF, tp)

	// list[T any] interface
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
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

	tp := typeDesc.NewNamed(`T`, ab.bakeAny())
	methods := map[string]typeDesc.TypeDesc{}

	// $len() int
	methods[`$len`] = ab.bakeIntFunc()

	// $recv() (T, bool)
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Return:     typeDesc.NewSolid(nil, ab.bakeReturnTuple(), tp),
	})
	getF = ab.proj.RegisterSignature(getF)
	methods[`$recv`] = typeDesc.NewSolid(nil, getF, tp)

	// $send(value T)
	setF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Params:     []typeDesc.Named{typeDesc.NewNamed(`value`, tp)},
	})
	setF = ab.proj.RegisterSignature(setF)
	methods[`$send`] = typeDesc.NewSolid(nil, setF, tp)

	// chan[T any] interface
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
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

	tpKey := typeDesc.NewNamed(`TKey`, ab.bakeAny())
	tpValue := typeDesc.NewNamed(`TValue`, ab.bakeAny())
	tp := []typeDesc.Named{tpKey, tpValue}
	methods := map[string]typeDesc.TypeDesc{}

	// $len() int
	methods[`$len`] = ab.bakeIntFunc()

	// $get(key TKey) (TValue, bool)
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: tp,
		Params:     []typeDesc.Named{typeDesc.NewNamed(`key`, tpKey)},
		Return:     typeDesc.NewSolid(nil, ab.bakeReturnTuple(), tpValue),
	})
	getF = ab.proj.RegisterSignature(getF)
	methods[`$get`] = typeDesc.NewSolid(nil, getF, tpKey, tpValue)

	// $set(key TKey, value TValue)
	setF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: tp,
		Params: []typeDesc.Named{
			typeDesc.NewNamed(`key`, tpKey),
			typeDesc.NewNamed(`value`, tpValue),
		},
	})
	setF = ab.proj.RegisterSignature(setF)
	methods[`$set`] = typeDesc.NewSolid(nil, setF, tpKey, tpValue)

	// map[TKey, TValue any] interface
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		TypeParams:   tp,
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
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

	tp := typeDesc.NewNamed(`T`, ab.bakeAny())
	methods := map[string]typeDesc.TypeDesc{}

	// $deref() T
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		TypeParams: []typeDesc.Named{tp},
		Return:     tp,
	})
	getF = ab.proj.RegisterSignature(getF)
	methods[`$deref`] = typeDesc.NewSolid(nil, getF, tp)

	// pointer[T any] interface
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		TypeParams:   []typeDesc.Named{tp},
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
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
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		Return: typeDesc.BasicFor[float32](ab.proj),
	})
	getF = ab.proj.RegisterSignature(getF)

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$real`] = getF // $real() float32
	methods[`$imag`] = getF // $imag() float32

	// complex64
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
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
	getF := typeDesc.NewSignature(typeDesc.SignatureArgs{
		Return: typeDesc.BasicFor[float64](ab.proj),
	})
	getF = ab.proj.RegisterSignature(getF)

	methods := map[string]typeDesc.TypeDesc{}
	methods[`$real`] = getF // $real() float32
	methods[`$imag`] = getF // $imag() float32

	// complex128
	t := typeDesc.NewInterface(typeDesc.InterfaceArgs{
		Methods:      methods,
		InitInherits: []typeDesc.Interface{ab.bakeAny()},
	})
	t = ab.proj.RegisterInterface(t)
	ab.baked[bakeKey] = t
	return t
}
