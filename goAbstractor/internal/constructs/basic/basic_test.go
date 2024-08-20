package basic

import (
	"go/types"
	"testing"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/Snow-Gremlin/goToolbox/testers/check"
)

func Test_Basic(t *testing.T) {
	tests := []struct {
		k   types.BasicKind
		exp types.BasicKind
		err string
	}{
		{k: types.Invalid, err: `may not use an invalid type in basic construct`},
		{k: types.Bool, exp: types.Bool},
		{k: types.Int, exp: types.Int},
		{k: types.Int8, exp: types.Int8},
		{k: types.Int16, exp: types.Int16},
		{k: types.Int32, exp: types.Int32},
		{k: types.Int64, exp: types.Int64},
		{k: types.Uint, exp: types.Uint},
		{k: types.Uint8, exp: types.Uint8},
		{k: types.Uint16, exp: types.Uint16},
		{k: types.Uint32, exp: types.Uint32},
		{k: types.Uint64, exp: types.Uint64},
		{k: types.Uintptr, exp: types.Uintptr},
		{k: types.Float32, exp: types.Float32},
		{k: types.Float64, exp: types.Float64},
		{k: types.Complex64, err: `unexpected complex type in basic construct`},
		{k: types.Complex128, err: `unexpected complex type in basic construct`},
		{k: types.String, exp: types.String},
		{k: types.UnsafePointer, exp: types.Uintptr},
		{k: types.UntypedBool, exp: types.Bool},
		{k: types.UntypedInt, exp: types.Int},
		{k: types.UntypedRune, exp: types.Int32},
		{k: types.UntypedFloat, exp: types.Float64},
		{k: types.UntypedComplex, err: `unexpected complex type in basic construct`},
		{k: types.UntypedString, exp: types.String},
		{k: types.UntypedNil, err: `unexpected untyped nil in basic construct`},
		{k: types.Byte, exp: types.Uint8},
		{k: types.Rune, exp: types.Int32},
	}

	for _, test := range tests {
		rt := types.Typ[test.k]
		if len(test.err) > 0 {
			check.MatchError(t, test.err).
				WithValue(`given`, rt.Name()).
				Panic(func() {
					newBasic(constructs.BasicArgs{RealType: rt})
				})
		} else {
			b := newBasic(constructs.BasicArgs{RealType: rt})
			rt2 := b.GoType().(*types.Basic)
			check.Equal(t, test.exp).
				WithValue(`given`, rt.Name()).
				WithValue(`gotten`, rt2.Name()).
				Assert(rt2.Kind())
		}
	}
}
