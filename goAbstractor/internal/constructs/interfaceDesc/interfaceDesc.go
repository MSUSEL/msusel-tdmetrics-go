package interfaceDesc

import (
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type interfaceDescImp struct {
	realType *types.Interface

	abstracts []constructs.Abstract
	exact     []constructs.TypeDesc
	approx    []constructs.TypeDesc

	inherits collections.SortedSet[constructs.InterfaceDesc]

	index int
}

func newInterfaceDesc(args constructs.InterfaceDescArgs) constructs.InterfaceDesc {
	assert.ArgNotNil(`real type`, args.RealType)
	return &interfaceDescImp{
		realType: args.RealType,

		abstracts: args.Abstracts,
		exact:     args.Exact,
		approx:    args.Approx,

		inherits: sortedSet.New[constructs.InterfaceDesc](),
	}
}

func (id *interfaceDescImp) IsTypeDesc()        {}
func (id *interfaceDescImp) IsInterfaceDesc()   {}
func (id *interfaceDescImp) Kind() kind.Kind    { return kind.InterfaceDesc }
func (id *interfaceDescImp) SetIndex(index int) { id.index = index }
func (id *interfaceDescImp) GoType() types.Type { return id.realType }

func (id *interfaceDescImp) IsUnion() bool {
	return len(id.approx)+len(id.exact) >= 2
}

func (id *interfaceDescImp) Implements(other constructs.InterfaceDesc) bool {
	return types.Implements(id.realType, other.(*interfaceDescImp).realType)
}

func (id *interfaceDescImp) AddInherits(it constructs.InterfaceDesc) constructs.InterfaceDesc {
	v, _ := id.inherits.TryAdd(it)
	return v
}

func (id *interfaceDescImp) Inherits() collections.SortedSet[constructs.InterfaceDesc] {
	return id.inherits
}

func (id *interfaceDescImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.InterfaceDesc](id, other, Comparer())
}

func Comparer() comp.Comparer[constructs.InterfaceDesc] {
	return func(a, b constructs.InterfaceDesc) int {
		aImp, bImp := a.(*interfaceDescImp), b.(*interfaceDescImp)
		return comp.Or(
			constructs.SliceComparerPend(aImp.abstracts, bImp.abstracts),
			constructs.SliceComparerPend(aImp.exact, bImp.exact),
			constructs.SliceComparerPend(aImp.approx, bImp.approx),
		)
	}
}

func (id *interfaceDescImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, id.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, id.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, id.index).
		AddNonZero(ctx2, `abstracts`, id.abstracts).
		AddNonZero(ctx2, `approx`, id.approx).
		AddNonZero(ctx2, `exact`, id.exact).
		AddNonZero(ctx2, `inherits`, id.inherits)
}
