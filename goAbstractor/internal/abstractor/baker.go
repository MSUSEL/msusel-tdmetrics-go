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

	t := typeDesc.NewInterface(types.NewInterfaceType(nil, nil))
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

	f := typeDesc.NewSignature(nil) // func() int
	f.SetReturn(typeDesc.BasicFor[int](ab.proj))
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

	t := typeDesc.NewStruct(nil)
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	t.AddField(`value`, tp, false)
	t.AddField(`ok`, typeDesc.BasicFor[bool](ab.proj), false)

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	t.AddFunc(`$len`, ab.bakeIntFunc()) // $len() int
	t.AddFunc(`$cap`, ab.bakeIntFunc()) // $cap() int

	getF := typeDesc.NewSignature(nil) // $get(index int) T
	getF.AppendTypeParam(tp)
	getF.AddParam(`index`, typeDesc.BasicFor[int](ab.proj))
	getF.SetReturn(tp)
	getF = ab.proj.RegisterSignature(getF)
	t.AddFunc(`$get`, typeDesc.NewSolid(nil, getF, tp))

	setF := typeDesc.NewSignature(nil) // $set(index int, value T)
	setF.AppendTypeParam(tp)
	setF.AddParam(`index`, typeDesc.BasicFor[int](ab.proj))
	setF.AddParam(`value`, tp)
	setF = ab.proj.RegisterSignature(setF)
	t.AddFunc(`$set`, typeDesc.NewSolid(nil, setF, tp))

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	t.AddFunc(`$len`, ab.bakeIntFunc()) // $len() int

	getF := typeDesc.NewSignature(nil) // $recv() (T, bool)
	getF.AppendTypeParam(tp)
	getF.SetReturn(typeDesc.NewSolid(nil, ab.bakeReturnTuple(), tp))
	getF = ab.proj.RegisterSignature(getF)
	t.AddFunc(`$recv`, typeDesc.NewSolid(nil, getF, tp))

	setF := typeDesc.NewSignature(nil) // $send(value T)
	setF.AppendTypeParam(tp)
	setF.AddParam(`value`, tp)
	setF = ab.proj.RegisterSignature(setF)
	t.AddFunc(`$send`, typeDesc.NewSolid(nil, setF, tp))

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())
	tpKey := t.AddTypeParam(`TKey`, ab.bakeAny())
	tpValue := t.AddTypeParam(`TValue`, ab.bakeAny())

	t.AddFunc(`$len`, ab.bakeIntFunc()) // $len() int

	getF := typeDesc.NewSignature(nil) // $get(key TKey) (TValue, bool)
	getF.AppendTypeParam(tpKey)
	getF.AppendTypeParam(tpValue)
	getF.AddParam(`key`, tpKey)
	getF.SetReturn(typeDesc.NewSolid(nil, ab.bakeReturnTuple(), tpValue))
	getF = ab.proj.RegisterSignature(getF)
	t.AddFunc(`$get`, typeDesc.NewSolid(nil, getF, tpKey, tpValue))

	setF := typeDesc.NewSignature(nil) // $set(key TKey, value TValue)
	getF.AppendTypeParam(tpKey)
	getF.AppendTypeParam(tpValue)
	setF.AddParam(`key`, tpKey)
	setF.AddParam(`value`, tpValue)
	setF = ab.proj.RegisterSignature(setF)
	t.AddFunc(`$set`, typeDesc.NewSolid(nil, setF, tpKey, tpValue))

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	getF := typeDesc.NewSignature(nil) // $deref() T
	getF.AppendTypeParam(tp)
	getF.SetReturn(tp)
	getF = ab.proj.RegisterSignature(getF)
	t.AddFunc(`$deref`, typeDesc.NewSolid(nil, getF, tp))

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())

	getF := typeDesc.NewSignature(nil) // func() float32
	getF.SetReturn(typeDesc.BasicFor[float32](ab.proj))
	getF = ab.proj.RegisterSignature(getF)

	t.AddFunc(`$real`, getF) // $real() float32
	t.AddFunc(`$imag`, getF) // $imag() float32

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

	t := typeDesc.NewInterface(nil)
	t.AppendInherits(ab.bakeAny())

	getF := typeDesc.NewSignature(nil) // func() float64
	getF.SetReturn(typeDesc.BasicFor[float64](ab.proj))
	getF = ab.proj.RegisterSignature(getF)

	t.AddFunc(`$real`, getF) // $real() float64
	t.AddFunc(`$imag`, getF) // $imag() float64

	t = ab.proj.RegisterInterface(t)
	ab.baked[bakeKey] = t
	return t
}
