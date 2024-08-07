package constructs

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Basic is a base type (e.g. bool, int, string, float64).
//
// This does not contain complex types, those are treated as an interface.
type Basic interface {
	TypeDesc
	_basic()
}

type BasicArgs struct {
	RealType *types.Basic

	// TypeName is only used if RealType is nil.
	TypeName string

	// Package must not be nil when RealType is nil.
	Package Package
}

type basicImp struct {
	realType *types.Basic
	typeName string
	index    int
}

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
		panic(terror.New(`unknown basic type name`).
			With(`name`, name))
	}
}

func newBasic(args BasicArgs) Basic {
	if !utils.IsNil(args.RealType) {
		return &basicImp{
			realType: args.RealType,
			typeName: normalizeBasicName(args.RealType.Name()),
		}
	}

	assert.ArgNotEmpty(`type name`, args.TypeName)
	assert.ArgNotNil(`package`, args.Package)

	typeName := normalizeBasicName(args.TypeName)
	pkg := args.Package.Source()
	tv, err := types.Eval(pkg.Fset, pkg.Types, token.NoPos, `(*`+typeName+`)(nil)`)
	if err != nil {
		panic(terror.New(`unable to create basic type from name`, err).
			With(`type name`, typeName))
	}
	realType := tv.Type.(*types.Pointer).Elem().(*types.Basic)

	return &basicImp{
		realType: realType,
		typeName: typeName,
	}
}

func (t *basicImp) _basic()            {}
func (t *basicImp) Kind() kind.Kind    { return kind.Basic }
func (t *basicImp) setIndex(index int) { t.index = index }
func (t *basicImp) GoType() types.Type { return t.realType }
func (t *basicImp) String() string     { return t.typeName }

func (t *basicImp) compareTo(other Construct) int {
	return strings.Compare(t.typeName, other.(*basicImp).typeName)
}

func (t *basicImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, t.index)
	}

	if ctx.IsKindShown() || ctx.IsIndexShown() {
		return jsonify.NewMap().
			AddIf(ctx, ctx.IsKindShown(), `kind`, t.Kind()).
			AddIf(ctx, ctx.IsIndexShown(), `index`, t.index).
			Add(ctx, `name`, t.typeName)
	}

	return jsonify.New(ctx, t.typeName)
}
