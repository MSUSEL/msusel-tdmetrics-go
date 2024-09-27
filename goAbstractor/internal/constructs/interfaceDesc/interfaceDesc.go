package interfaceDesc

import (
	"go/token"
	"go/types"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type interfaceDescImp struct {
	hint     hint.Hint
	realType types.Type

	pinnedPkg constructs.Package
	abstracts []constructs.Abstract
	exact     []constructs.TypeDesc
	approx    []constructs.TypeDesc

	inherits collections.SortedSet[constructs.InterfaceDesc]

	index int
	alive bool
}

func newInterfaceDesc(args constructs.InterfaceDescArgs) constructs.InterfaceDesc {
	assert.ArgHasNoNils(`abstracts`, args.Abstracts)
	assert.ArgHasNoNils(`exact`, args.Exact)
	assert.ArgHasNoNils(`approx`, args.Approx)

	if utils.IsNil(args.RealType) {
		assert.ArgNotNil(`package`, args.Package)

		switch args.Hint {
		case hint.Pointer:
			// TODO: Implement

		case hint.List:
			// TODO: Implement

		case hint.Map:
			// TODO: Implement

		case hint.Chan:
			// TODO: Implement
			//args.RealType =

		case hint.Complex64:
			args.RealType = types.Typ[types.Complex64]

		case hint.Complex128:
			args.RealType = types.Typ[types.Complex128]

		default: // hint.None, hint.Comparable
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
	}
	assert.ArgNotNil(`real type`, args.RealType)

	return &interfaceDescImp{
		hint:     args.Hint,
		realType: args.RealType,

		pinnedPkg: args.PinnedPkg,
		abstracts: args.Abstracts,
		exact:     args.Exact,
		approx:    args.Approx,

		inherits: sortedSet.New(Comparer()),
	}
}

func (id *interfaceDescImp) IsTypeDesc()      {}
func (id *interfaceDescImp) IsInterfaceDesc() {}

func (id *interfaceDescImp) Kind() kind.Kind     { return kind.InterfaceDesc }
func (id *interfaceDescImp) Index() int          { return id.index }
func (id *interfaceDescImp) SetIndex(index int)  { id.index = index }
func (id *interfaceDescImp) Alive() bool         { return id.alive }
func (id *interfaceDescImp) SetAlive(alive bool) { id.alive = alive }
func (id *interfaceDescImp) Hint() hint.Hint     { return id.hint }
func (id *interfaceDescImp) GoType() types.Type  { return id.realType }

func (id *interfaceDescImp) Abstracts() []constructs.Abstract  { return id.abstracts }
func (id *interfaceDescImp) Exact() []constructs.TypeDesc      { return id.exact }
func (id *interfaceDescImp) Approx() []constructs.TypeDesc     { return id.approx }
func (id *interfaceDescImp) PinnedPackage() constructs.Package { return id.pinnedPkg }

func (id *interfaceDescImp) IsPinned() bool {
	return utils.IsNil(id.pinnedPkg)
}

func (id *interfaceDescImp) IsGeneral() bool {
	return len(id.approx)+len(id.exact) >= 2
}

func (id *interfaceDescImp) Implements(other constructs.InterfaceDesc) bool {
	rtIt, ok := other.(*interfaceDescImp).realType.(*types.Interface)
	return ok && types.Implements(id.realType, rtIt)
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
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			constructs.ComparerPend(aImp.pinnedPkg, bImp.pinnedPkg),
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
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, id.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, id.Kind(), id.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, id.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, id.index).
		AddNonZero(ctx.OnlyIndex(), `package`, id.pinnedPkg).
		AddNonZero(ctx.OnlyIndex(), `abstracts`, id.abstracts).
		AddNonZero(ctx.Short(), `approx`, id.approx).
		AddNonZero(ctx.Short(), `exact`, id.exact).
		AddNonZero(ctx.OnlyIndex(), `inherits`, id.inherits.ToSlice())
}

func (id *interfaceDescImp) String() string {
	internals := ``
	if len(id.abstracts) > 0 {
		internals += enumerator.Enumerate(id.abstracts...).Join(`; `) + `; `
	}
	if len(id.exact) > 0 {
		internals += enumerator.Enumerate(id.exact...).Join(`|`)
		if len(id.approx) > 0 {
			internals += `|~` + enumerator.Enumerate(id.approx...).Join(`|~`) + `; `
		} else {
			internals += `; `
		}
	} else if len(id.approx) > 0 {
		internals += `~` + enumerator.Enumerate(id.approx...).Join(`|~`) + `; `
	}
	if len(internals) <= 0 {
		return `any`
	}
	head := ``
	if id.IsPinned() {
		head = id.pinnedPkg.Path() + `:`
	}
	return head + `interface{ ` + internals + `}`
}
