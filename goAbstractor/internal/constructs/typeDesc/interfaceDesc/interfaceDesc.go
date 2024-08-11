package interfaceDesc

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/typeDesc/signature"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
)

const Kind = `interfaceDesc`

// InterfaceDesc is a named interface typically explicitly defined at the given
// location in the source code. The underlying type description
// can be a class or interface with optional parameter types.
//
// If type parameters are given then the interface is generic.
// Instances with realized versions of the interface,
// are added for each used instance in the source code. If there
// are no instances then the generic interface isn't used.
type InterfaceDesc interface {
	typeDesc.TypeDesc
	_interfaceDesc()

	// IsUnion indicates if there is two or more exact or approximate types.
	IsUnion() bool

	AddInherits(it InterfaceDesc) InterfaceDesc

	Inherits() collections.ReadonlySortedSet[InterfaceDesc]
}

type Args struct {
	RealType *types.Interface

	// Methods is the set of signatures for this interface.
	Signatures []signature.Signature

	// Exact types are like `string|int|bool` where the
	// data type must match exactly.
	Exact []typeDesc.TypeDesc

	// Approx types are like `~string|~int` where the data type
	// may be exact or an extension of the base type.
	Approx []typeDesc.TypeDesc
}

type interfaceDescImp struct {
	realType *types.Interface

	signatures []signature.Signature
	exact      []typeDesc.TypeDesc
	approx     []typeDesc.TypeDesc

	inherits collections.SortedSet[InterfaceDesc]

	index int
}

func newInterfaceDesc(args Args) InterfaceDesc {
	assert.ArgNotNil(`real type`, args.RealType)
	return &interfaceDescImp{
		realType: args.RealType,

		signatures: args.Signatures,
		exact:      args.Exact,
		approx:     args.Approx,

		inherits: sortedSet.New[InterfaceDesc](),
	}
}

func (id *interfaceDescImp) _interfaceDesc()    {}
func (id *interfaceDescImp) Kind() string       { return Kind }
func (id *interfaceDescImp) SetIndex(index int) { id.index = index }
func (id *interfaceDescImp) GoType() types.Type { return id.realType }

func (id *interfaceDescImp) IsUnion() bool {
	return len(id.approx)+len(id.exact) >= 2
}

func (id *interfaceDescImp) AddInherits(it InterfaceDesc) InterfaceDesc {
	v, _ := id.inherits.TryAdd(it)
	return v
}

func (id *interfaceDescImp) Inherits() collections.ReadonlySortedSet[InterfaceDesc] {
	return id.inherits.Readonly()
}

func (id *interfaceDescImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[InterfaceDesc](id, other)
}

func Comparer() comp.Comparer[InterfaceDesc] {
	return func(a, b InterfaceDesc) int {
		aImp, bImp := a.(*interfaceDescImp), b.(*interfaceDescImp)
		return comp.Or(
			constructs.SliceComparerPend(aImp.signatures, bImp.signatures),
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
		AddIf(ctx, ctx.IsKindShown(), `kind`, Kind).
		AddIf(ctx, ctx.IsIndexShown(), `index`, id.index).
		AddNonZero(ctx2, `signatures`, id.signatures).
		AddNonZero(ctx2, `approx`, id.approx).
		AddNonZero(ctx2, `exact`, id.exact).
		AddNonZero(ctx2, `inherits`, id.inherits)
}
