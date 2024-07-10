package constructs

import (
	"errors"
	"fmt"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

// Solid represents a generic type that has been resolved to a specific type
// with specific type parameters, e.g. List<T> might be resolved to List<int>.
// The type parameter resolution may be referencing another type parameter,
// e.g. a method signature inside a generic interface.
type Solid interface {
	TypeDesc
	_solid()
}

func newSolid(realType types.Type, target TypeDesc, tp ...TypeDesc) Solid {
	if len(tp) <= 0 {
		panic(fmt.Sprintf(`a solid type requires at least one type parameter %v`, target))
	}
	if utils.IsNil(realType) {
		panic(errors.New(`must provide a real type for a solid`))
	}
	return &solidImp{
		realType:   realType,
		target:     target,
		typeParams: tp,
	}
}

type solidImp struct {
	realType types.Type

	index      int
	target     TypeDesc
	typeParams []TypeDesc
}

func (ts *solidImp) _solid() {}

func (ts *solidImp) Visit(v Visitor) {
	visitTest(v, ts.target)
	visitList(v, ts.typeParams)
}

func (ts *solidImp) SetIndex(index int) {
	ts.index = index
}

func (ts *solidImp) GoType() types.Type {
	return ts.realType
}

func (ts *solidImp) Equal(other TypeDesc) bool {
	return equalTest(ts, other, func(a, b *solidImp) bool {
		return equal(a.target, b.target) &&
			equalList(a.typeParams, b.typeParams)
	})
}

func (ts *solidImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, ts.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, `solid`).
		AddNonZero(ctx2, `target`, ts.target).
		AddNonZero(ctx2, `typeParams`, ts.typeParams)
}

func (ts *solidImp) String() string {
	return jsonify.ToString(ts)
}
