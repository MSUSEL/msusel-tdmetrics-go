package baker

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

// MethodByName gets a method by name or returns nil.
//
// The given args are hints at what the types of the builtin
// methods params are. The first argument is the return type
// of the method. Any arg in the args may be nil.
func (b *bakerImp) MethodByName(name string, args []types.Type) constructs.Method {
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

// BakeAppend creates the builtin append function.
//
//	func append(slice []Type, elems ...Type) []Type
func (b *bakerImp) BakeAppend(args []types.Type) constructs.Method {
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
func (b *bakerImp) BakeClear(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeClose creates the builtin close function.
//
//	func close(c chan<- Type)
func (b *bakerImp) BakeClose(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeComplex creates the builtin complex function.
//
//	func complex(r, i FloatType) ComplexType
func (b *bakerImp) BakeComplex(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeCopy creates the builtin copy function.
//
//	func copy(dst, src []Type) int
func (b *bakerImp) BakeCopy(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeDelete creates the builtin delete function.
//
//	func delete(m map[Type]Type1, key Type)
func (b *bakerImp) BakeDelete(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeImag creates the builtin imag function.
//
//	func imag(c ComplexType) FloatType
func (b *bakerImp) BakeImag(args []types.Type) constructs.Method {
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
func (b *bakerImp) BakeMake(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeMax creates the builtin max function.
//
//	func max[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMax(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeMin creates the builtin min function.
//
//	func min[T cmp.Ordered](x T, y ...T) T
func (b *bakerImp) BakeMin(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
}

// BakeNew creates the builtin new function.
//
//	func new(Type) *Type
func (b *bakerImp) BakeNew(args []types.Type) constructs.Method {
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

		// args any
		ret := b.proj.NewArgument(constructs.ArgumentArgs{
			Name: `args`,
			Type: b.BakeAny(),
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
			Name:      `println`,
			Exported:  true,
			Signature: sig,
		})
	})
}

// BakeReal creates the builtin real function.
//
//	func real(c ComplexType) FloatType
func (b *bakerImp) BakeReal(args []types.Type) constructs.Method {
	assert.NotImplemented()
	return nil // TODO: Implement
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
		})
	})
}
