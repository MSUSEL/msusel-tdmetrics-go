package baker

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/abstractor/instantiator"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

// MethodByName gets a method by name or returns nil.
// This may return a method or method instance.
//
// The given args are hints at what the types of the builtin
// methods params are. The first argument is the return type
// of the method. Any arg in the args may be nil.
func (b *bakerImp) MethodByName(name string, args []constructs.TypeDesc) constructs.Construct {
	switch name {
	case `append`:
		return b.BakeAppend(args)
	case `cap`:
		return b.BakeCap()
	case `clear`:
		return b.BakeClear(args)
	case `close`:
		return b.BakeClose(args)
	case `complex`:
		return b.BakeComplex(args)
	case `copy`:
		return b.BakeCopy(args)
	case `delete`:
		return b.BakeDelete(args)
	case `imag`:
		return b.BakeImag(args)
	case `len`:
		return b.BakeLen()
	case `make`:
		return b.BakeMake(args)
	case `max`:
		return b.BakeMax(args)
	case `min`:
		return b.BakeMin(args)
	case `new`:
		return b.BakeNew(args)
	case `panic`:
		return b.BakePanic()
	case `print`:
		return b.BakePrint()
	case `println`:
		return b.BakePrintln()
	case `real`:
		return b.BakeReal(args)
	case `recover`:
		return b.BakeRecover()
	default:
		return nil
	}
}

// TODO: Remove once all of the methods are implemented.
func notImplementedMethod(name string, args []constructs.TypeDesc) {
	text := name + ` => [`
	for i, arg := range args {
		if i > 0 {
			text += `, `
		}
		text += arg.String()
	}
	text += `]`

	fmt.Println(text)
	assert.NotImplemented()
}

// BakeAppend creates the builtin append function.
//
//	func append(slice []Type, elems ...Type) []Type
func (b *bakerImp) BakeAppend(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`append`, args)
	return nil // TODO: Implement
}

// BakeCap creates the builtin cap function.
//
//	func cap(v Type) int
func (b *bakerImp) BakeCap() constructs.Method {
	return b.bakeLenCap(`cap`)
}

// BakeLen creates the builtin len function.
//
//	func len(v Type) int
func (b *bakerImp) BakeLen() constructs.Method {
	return b.bakeLenCap(`len`)
}

func (b *bakerImp) bakeLenCap(name string) constructs.Method {
	return bakeOnce(b, name, func() constructs.Method {
		pkg := b.BakeBuiltin()

		// v any
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `v`,
			Type: b.BakeAny(),
		})

		// <unnamed> int
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.bakeBasic(types.Int),
		})

		// func(v any) int
		sig := b.proj.NewSignature(constructs.SignatureArgs{
			Params:  []constructs.Argument{param},
			Results: []constructs.Argument{ret},
			Package: pkg.Source(),
		})

		// func <name>(v any) int
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:   pkg,
			Name:      name,
			Exported:  true,
			Signature: sig,
			Location:  locs.NoLoc(),
		})
	})
}

// BakeClear creates the builtin clear function.
//
//	func clear[T ~[]Type | ~map[Type]Type1](t T)
func (b *bakerImp) BakeClear(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`clear`, args)
	return nil // TODO: Implement
}

// BakeClose creates the builtin close function.
//
//	func close(c chan<- Type)
func (b *bakerImp) BakeClose(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`close`, args)
	return nil // TODO: Implement
}

// BakeCopy creates the builtin copy function.
//
//	func copy(dst, src []Type) int
func (b *bakerImp) BakeCopy(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`copy`, args)
	return nil // TODO: Implement
}

// BakeDelete creates the builtin delete function.
//
//	func delete(m map[Type]Type1, key Type)
func (b *bakerImp) BakeDelete(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`delete`, args)
	return nil // TODO: Implement
}

// BakeComplex creates the builtin complex function.
//
//	func complex(r, i FloatType) ComplexType
func (b *bakerImp) BakeComplex(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`complex`, args)
	return nil // TODO: Implement
}

// BakeImag creates the builtin imag function.
//
//	func imag(c ComplexType) FloatType
func (b *bakerImp) BakeImag(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`imag`, args)
	return nil // TODO: Implement
}

// BakeReal creates the builtin real function.
//
//	func real(c ComplexType) FloatType
func (b *bakerImp) BakeReal(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`real`, args)
	return nil // TODO: Implement
}

// BakeMake creates the builtin make function.
//
//	func make(t Type, size ...IntegerType) Type
func (b *bakerImp) BakeMake(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`make`, args)
	return nil // TODO: Implement
}

// BakeNew creates the builtin new function.
//
//	func new(Type) *Type
func (b *bakerImp) BakeNew(args []constructs.TypeDesc) constructs.Method {
	notImplementedMethod(`new`, args)
	return nil // TODO: Implement
}

// BakeMax creates the builtin max function.
//
//	func max[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMax(args []constructs.TypeDesc) constructs.Construct {
	return b.bakeMinMax(`max`, args)
}

// BakeMin creates the builtin min function.
//
//	func min[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMin(args []constructs.TypeDesc) constructs.Construct {
	return b.bakeMinMax(`min`, args)
}

func (b *bakerImp) bakeMinMax(name string, args []constructs.TypeDesc) constructs.Construct {
	assert.ArgNotEmpty(`args`, args)
	instTyp := args[0] // use the result type to determine the method
	assert.ArgNotNil(`typ arg`, instTyp)

	return bakeOnce(b, name+` `+instTyp.String(), func() constructs.Construct {
		// func <name>[T any](v ...T) T
		gen := b.bakeMinMaxGeneric(name)
		return instantiator.Method(b.proj, gen, instTyp)
	})
}

func (b *bakerImp) bakeMinMaxGeneric(name string) constructs.Method {
	return bakeOnce(b, name, func() constructs.Method {
		pkg := b.BakeBuiltin()

		// T any
		tp := b.proj.NewTypeParam(constructs.TypeParamArgs{
			Name: `T`,
			Type: b.BakeAny(),
		})

		// List[T]
		listT := instantiator.InterfaceDecl(b.proj, nil, b.BakeList(), tp)

		// x []T
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `x`,
			Type: listT,
		})

		// <unnamed> T
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: tp,
		})

		tpRt := types.NewTypeParam(types.NewTypeName(token.NoPos, pkg.Source().Types, `T`, nil), types.NewInterfaceType(nil, nil))
		rt := types.NewSignatureType(nil, nil, []*types.TypeParam{tpRt},
			types.NewTuple(types.NewParam(token.NoPos, pkg.Source().Types, `x`, tpRt)),
			types.NewTuple(types.NewParam(token.NoPos, pkg.Source().Types, ``, tpRt)), true)

		// func(x ...T) T
		sig := b.proj.NewSignature(constructs.SignatureArgs{
			RealType: rt,
			Params:   []constructs.Argument{param},
			Variadic: true,
			Results:  []constructs.Argument{ret},
			Package:  pkg.Source(),
		})

		// func <name>[T any](v ...T) T
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:    pkg,
			Name:       name,
			Exported:   true,
			TypeParams: []constructs.TypeParam{tp},
			Signature:  sig,
			Location:   locs.NoLoc(),
		})
	})
}

// BakePanic creates the builtin panic function.
//
//	func panic(v any)
func (b *bakerImp) BakePanic() constructs.Method {
	return bakeOnce(b, `panic`, func() constructs.Method {
		pkg := b.BakeBuiltin()

		// v any
		param := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `v`,
			Type: b.BakeAny(),
		})

		// func(v any)
		sig := b.proj.NewSignature(constructs.SignatureArgs{
			Params:  []constructs.Argument{param},
			Package: pkg.Source(),
		})

		// func panic(v any)
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:   pkg,
			Name:      `panic`,
			Exported:  true,
			Signature: sig,
			Location:  locs.NoLoc(),
		})
	})
}

// BakePrint creates the builtin print function.
//
//	func print(args ...Type)
func (b *bakerImp) BakePrint() constructs.Method {
	return b.bakePrintFunc(`print`)
}

// BakePrintln creates the builtin println function.
//
//	func println(args ...Type)
func (b *bakerImp) BakePrintln() constructs.Method {
	return b.bakePrintFunc(`println`)
}

func (b *bakerImp) bakePrintFunc(name string) constructs.Method {
	return bakeOnce(b, name, func() constructs.Method {
		pkg := b.BakeBuiltin()

		// List[any]
		anyList := instantiator.InterfaceDecl(b.proj, nil, b.BakeList(), b.BakeAny())

		// args []any
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `args`,
			Type: anyList,
		})

		// func(args ...any)
		sig := b.proj.NewSignature(constructs.SignatureArgs{
			Params:   []constructs.Argument{ret},
			Variadic: true,
			Package:  pkg.Source(),
		})

		// func println(args ...any)
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:   pkg,
			Name:      name,
			Exported:  true,
			Signature: sig,
			Location:  locs.NoLoc(),
		})
	})
}

// BakeRecover creates the builtin recover function.
//
//	func recover() any
func (b *bakerImp) BakeRecover() constructs.Method {
	return bakeOnce(b, `recover`, func() constructs.Method {
		pkg := b.BakeBuiltin()

		// <unnamed> any
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Type: b.BakeAny(),
		})

		// func() any
		sig := b.proj.NewSignature(constructs.SignatureArgs{
			Results: []constructs.Argument{ret},
			Package: pkg.Source(),
		})

		// func recover() any
		return b.proj.NewMethod(constructs.MethodArgs{
			Package:   pkg,
			Name:      `recover`,
			Exported:  true,
			Signature: sig,
			Location:  locs.NoLoc(),
		})
	})
}
