package interfaceDesc

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

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

	id any
}

func newInterfaceDesc(args constructs.InterfaceDescArgs) constructs.InterfaceDesc {
	assert.ArgHasNoNils(`abstracts`, args.Abstracts)
	assert.ArgHasNoNils(`exact`, args.Exact)
	assert.ArgHasNoNils(`approx`, args.Approx)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		methods := make([]*types.Func, len(args.Abstracts))
		for i, abstract := range args.Abstracts {
			methods[i] = types.NewFunc(token.NoPos, args.Package.Types,
				abstract.Name(), abstract.Signature().GoType().(*types.Signature))
		}

		embedded := []types.Type{}
		if count1, count2 := len(args.Exact), len(args.Approx); count1+count2 > 0 {
			terms := make([]*types.Term, count1+count2)
			for i, exact := range args.Exact {
				terms[i] = types.NewTerm(false, exact.GoType())
			}
			for i, approx := range args.Approx {
				terms[i+count1] = types.NewTerm(true, approx.GoType())
			}
			embedded = append(embedded, types.NewUnion(terms))
		}

		args.RealType = types.NewInterfaceType(methods, embedded).Complete()
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &interfaceDescImp{
		realType: args.RealType,

		abstracts: args.Abstracts,
		exact:     args.Exact,
		approx:    args.Approx,

		inherits: sortedSet.New(Comparer()),
	}
}

func (id *interfaceDescImp) IsTypeDesc()      {}
func (id *interfaceDescImp) IsInterfaceDesc() {}

func (id *interfaceDescImp) Kind() kind.Kind    { return kind.InterfaceDesc }
func (id *interfaceDescImp) Id() any            { return id.id }
func (id *interfaceDescImp) SetId(ident any)    { id.id = ident }
func (id *interfaceDescImp) GoType() types.Type { return id.realType }

func (id *interfaceDescImp) Abstracts() []constructs.Abstract { return id.abstracts }
func (id *interfaceDescImp) Exact() []constructs.TypeDesc     { return id.exact }
func (id *interfaceDescImp) Approx() []constructs.TypeDesc    { return id.approx }

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

func (id *interfaceDescImp) RemoveTempReferences() {
	for i, ap := range id.approx {
		id.approx[i] = constructs.ResolvedTempReference(ap)
	}
	for i, ex := range id.exact {
		id.exact[i] = constructs.ResolvedTempReference(ex)
	}
}

func (id *interfaceDescImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, id.id)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, id.Kind()).
		AddIf(ctx, ctx.IsIdShown(), `id`, id.id).
		AddNonZero(ctx2, `abstracts`, id.abstracts).
		AddNonZero(ctx2, `approx`, id.approx).
		AddNonZero(ctx2, `exact`, id.exact).
		AddNonZero(ctx2, `inherits`, id.inherits.ToSlice())
}

func (id *interfaceDescImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(`interface{ `)
	if len(id.abstracts) > 0 {
		buf.WriteString(enumerator.Enumerate(id.abstracts...).Join(`; `) + `; `)
	}
	if len(id.exact) > 0 {
		buf.WriteString(enumerator.Enumerate(id.exact...).Join(`|`) + `; `)
	}
	if len(id.approx) > 0 {
		buf.WriteString(`~` + enumerator.Enumerate(id.approx...).Join(`|~`) + `; `)
	}
	buf.WriteString(`}`)
	return buf.String()
}
