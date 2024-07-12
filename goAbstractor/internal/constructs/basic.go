package constructs

import (
	"errors"
	"fmt"
	"go/token"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/utils"
	"golang.org/x/tools/go/packages"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Basic interface {
		TypeDesc
		_basic()
	}

	basicImp struct {
		typ   *types.Basic
		name  string
		index int
	}
)

func normalizeBasicName(name string) string {
	name, _ = strings.CutPrefix(name, `untyped `)
	switch name {
	case `byte`:
		return `uint8`
	case `rune`:
		return `int32`
	case `float`:
		return `float64`
	case `Pointer`:
		return `uintptr`
	case `int`, `uint`, `int8`, `uint8`, `int16`, `uint16`, `int32`, `uint32`,
		`int64`, `uint64`, `float32`, `float64`, `string`, `bool`, `uintptr`:
		return name
	default:
		panic(fmt.Errorf(`unexpected basic type: %q`, name))
	}
}

func newBasic(typ *types.Basic) Basic {
	if utils.IsNil(typ) {
		panic(errors.New(`may not create a new basic with a nil type`))
	}
	return &basicImp{
		typ:  typ,
		name: normalizeBasicName(typ.Name()),
	}
}

func newBasicFromName(pkg *packages.Package, typeName string) Basic {
	typeName = normalizeBasicName(typeName)
	tv, err := types.Eval(pkg.Fset, pkg.Types, token.NoPos, `(*`+typeName+`)(nil)`)
	if err != nil {
		panic(fmt.Errorf(`unable to create basic type of %s: %w`, typeName, err))
	}
	typ := tv.Type.(*types.Pointer).Elem().(*types.Basic)
	return newBasic(typ)
}

func (t *basicImp) _basic() {}

func (t *basicImp) Visit(v Visitor) {}

func (t *basicImp) SetIndex(index int) {
	t.index = index
}

func (t *basicImp) GoType() types.Type {
	return t.typ
}

func (t *basicImp) Equal(other Construct) bool {
	return equalTest(t, other, func(a, b *basicImp) bool {
		return a.name == b.name
	})
}

func (t *basicImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	if ctx.IsKindShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, `basic`).
			Add(ctx, `name`, t.name)
	}

	return jsonify.New(ctx, t.name)
}

func (t *basicImp) String() string {
	return jsonify.ToString(t)
}
