package abstractor

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"

// bakeAny bakes in an interface to represent "any"
// the base object that (almost) all other types inherit from.
func (ab *abstractor) bakeAny() *typeDesc.Interface {
	const bakeKey = `any`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Interface)
	}

	t := typeDesc.NewInterface()
	t = ab.registerInterface(t)
	ab.baked[bakeKey] = t
	return t
}

func (ab *abstractor) bakeIntFunc() *typeDesc.Signature {
	const bakeKey = `func() int`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Signature)
	}

	f := typeDesc.NewSignature() // $len() int
	f.Return = typeDesc.NewBasic(`int`)
	f = ab.registerSignature(f)
	ab.baked[bakeKey] = f
	return f
}

// bakeArray bakes in an interface to represent a Go array:
//
//	type array[T any] interface {
//		$len() int
//		$get(index int) T
//		$set(index int, value T)
//	}
//
// Note: Doesn't currently have cap as defined in reflect.
func (ab *abstractor) bakeArray() *typeDesc.Interface {
	const bakeKey = `array[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Interface)
	}

	t := typeDesc.NewInterface()
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	lenF := typeDesc.NewSignature() // $len() int
	lenF.Return = typeDesc.NewBasic(`int`)
	t.AddFunc(`$len`, ab.registerSignature(lenF))

	getF := typeDesc.NewSignature() // $get(index int) T
	getF.AppendTypeParam(tp)
	getF.AddParam(`index`, typeDesc.NewBasic(`int`))
	getF.Return = tp
	t.AddFunc(`$get`, typeDesc.NewSolid(ab.registerSignature(getF), tp))

	setF := typeDesc.NewSignature() // $set(index int, value T)
	setF.AppendTypeParam(tp)
	setF.AddParam(`index`, typeDesc.NewBasic(`int`))
	setF.AddParam(`value`, tp)
	t.AddFunc(`$set`, typeDesc.NewSolid(ab.registerSignature(setF), tp))

	ab.proj.AllInterfaces = append(ab.proj.AllInterfaces, t)
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
func (ab *abstractor) bakeChan() *typeDesc.Interface {
	const bakeKey = `chan[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Interface)
	}

	t := typeDesc.NewInterface()
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	lenF := typeDesc.NewSignature() // $len() int
	lenF.Return = typeDesc.NewBasic(`int`)
	t.AddFunc(`$len`, ab.registerSignature(lenF))

	getFR := typeDesc.NewStruct()
	getFR.AppendTypeParam(tp)
	getFR.AddField(`value`, tp, false)
	getFR.AddField(`ok`, typeDesc.NewBasic(`bool`), false)

	getF := typeDesc.NewSignature() // $recv() (T, bool)
	getF.AppendTypeParam(tp)
	getF.Return = typeDesc.NewSolid(ab.registerStruct(getFR), tp)
	t.AddFunc(`$recv`, typeDesc.NewSolid(ab.registerSignature(getF), tp))

	setF := typeDesc.NewSignature() // $send(value T)
	setF.AppendTypeParam(tp)
	setF.AddParam(`value`, tp)
	t.AddFunc(`$send`, typeDesc.NewSolid(ab.registerSignature(setF), tp))

	ab.proj.AllInterfaces = append(ab.proj.AllInterfaces, t)
	ab.baked[bakeKey] = t
	return t
}

// This bakes in an interface to represent a Go map:
//
//	type map[TKey, TValue any] interface {
//		$len() int
//		$get(key TKey) (TValue, bool)
//		$set(key TKey, value TValue)
//	}
//
// Note: Doesn't currently require Key to be comparable as defined in reflect.
func (ab *abstractor) bakeMap() *typeDesc.Interface {
	const bakeKey = `map[TKey, TValue any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Interface)
	}

	t := typeDesc.NewInterface()
	tpKey := t.AddTypeParam(`TKey`, ab.bakeAny())
	tpVal := t.AddTypeParam(`TValue`, ab.bakeAny())

	lenF := typeDesc.NewSignature() // $len() int
	lenF.Return = typeDesc.NewBasic(`int`)
	t.AddFunc(`$len`, ab.registerSignature(lenF))

	getFR := typeDesc.NewStruct()
	getFR.AppendTypeParam(tpVal)
	getFR.AddField(`value`, tpVal, false)
	getFR.AddField(`ok`, typeDesc.NewBasic(`bool`), false)

	getF := typeDesc.NewSignature() // $get(key TKey) (TValue, bool)
	getF.AppendTypeParam(tpKey)
	getF.AppendTypeParam(tpVal)
	getF.AddParam(`key`, tpKey)
	getF.Return = getFR
	t.AddFunc(`$get`, typeDesc.NewSolid(ab.registerSignature(getF), tpKey, tpVal))

	setF := typeDesc.NewSignature() // $set(key TKey, value TValue)
	getF.AppendTypeParam(tpKey)
	getF.AppendTypeParam(tpVal)
	setF.AddParam(`key`, tpKey)
	setF.AddParam(`value`, tpVal)
	t.AddFunc(`$set`, typeDesc.NewSolid(ab.registerSignature(setF), tpKey, tpVal))

	ab.proj.AllInterfaces = append(ab.proj.AllInterfaces, t)
	ab.baked[bakeKey] = t
	return t
}

func (ab *abstractor) bakePointer() *typeDesc.Interface {
}

// bakeSlice bakes in an interface to represent a Go array:
//
//	type slice[T any] interface {
//		$len() int
//		$cap() int
//		$get(index int) T
//		$set(index int, value T)
//	}
func (ab *abstractor) bakeSlice() *typeDesc.Interface {
	const bakeKey = `array[T any]`
	if t, has := ab.baked[bakeKey]; has {
		return t.(*typeDesc.Interface)
	}

	t := typeDesc.NewInterface()
	tp := t.AddTypeParam(`T`, ab.bakeAny())

	lenF := typeDesc.NewSignature()
	lenF.Return = typeDesc.NewBasic(`int`)
	lenF = ab.registerSignature(lenF)
	t.AddFunc(`$len`, lenF) // $len() int
	t.AddFunc(`$cap`, lenF) // $cap() int

	getF := typeDesc.NewSignature() // $get(index int) T
	getF.AppendTypeParam(tp)
	getF.AddParam(`index`, typeDesc.NewBasic(`int`))
	getF.Return = tp
	t.AddFunc(`$get`, typeDesc.NewSolid(ab.registerSignature(getF), tp))

	setF := typeDesc.NewSignature() // $set(index int, value T)
	setF.AppendTypeParam(tp)
	setF.AddParam(`index`, typeDesc.NewBasic(`int`))
	setF.AddParam(`value`, tp)
	t.AddFunc(`$set`, typeDesc.NewSolid(ab.registerSignature(setF), tp))

	ab.proj.AllInterfaces = append(ab.proj.AllInterfaces, t)
	ab.baked[bakeKey] = t
	return t
}
