package interfaceDesc

import (
	"go/token"
	"go/types"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/hint"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/innate"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type interfaceDescImp struct {
	hint     hint.Hint
	realType types.Type

	pinnedPkg constructs.Package
	abstracts []constructs.Abstract
	exact     []constructs.TypeDesc
	approx    []constructs.TypeDesc
	additions []constructs.Abstract

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
			// Get $deref's only resulting real type.
			ret := constructs.FindSigByName(args.Abstracts, innate.Deref).Results()[0].Type().GoType()
			args.RealType = types.NewPointer(ret)

		case hint.List:
			// Get $get's only resulting real type.
			ret := constructs.FindSigByName(args.Abstracts, innate.Get).Results()[0].Type().GoType()
			args.RealType = types.NewSlice(ret)

		case hint.Map:
			// Get $set's two parameter real types.
			params := constructs.FindSigByName(args.Abstracts, innate.Set).Params()
			keyRet := params[0].Type().GoType()
			valRet := params[1].Type().GoType()
			args.RealType = types.NewMap(keyRet, valRet)

		case hint.Chan:
			// Get $send's only parameter real type.
			ret := constructs.FindSigByName(args.Abstracts, innate.Send).Params()[0].Type().GoType()
			args.RealType = types.NewSlice(ret)

		case hint.Complex64:
			args.RealType = types.Typ[types.Complex64]

		case hint.Complex128:
			args.RealType = types.Typ[types.Complex128]

		default: // hint.None, hint.Comparable
			methods := make([]*types.Func, 0, len(args.Abstracts))
			for _, abstract := range args.Abstracts {
				if name := abstract.Name(); !innate.Is(name) {
					sig := abstract.Signature().GoType().(*types.Signature)
					m := types.NewFunc(token.NoPos, args.Package.Types, name, sig)
					methods = append(methods, m)
				}
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

			if args.Hint == hint.Comparable && len(methods) <= 0 && len(embedded) <= 0 {
				args.RealType = types.Universe.Lookup("comparable").Type().Underlying()
			} else {
				args.RealType = types.NewInterfaceType(methods, embedded).Complete()
			}
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
	return !utils.IsNil(id.pinnedPkg)
}

func (id *interfaceDescImp) IsGeneral() bool {
	return len(id.approx)+len(id.exact) >= 2
}

func (id *interfaceDescImp) Implements(other constructs.InterfaceDesc) bool {
	thisIt := id.realType
	otherIt, ok := other.GoType().(*types.Interface)
	return ok && types.Implements(thisIt, otherIt)
}

func (id *interfaceDescImp) AdditionalAbstracts() []constructs.Abstract {
	return id.additions
}

func (id *interfaceDescImp) SetAdditionalAbstracts(abstracts []constructs.Abstract) {
	id.additions = abstracts
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

func (id *interfaceDescImp) RemoveTempReferences(required bool) {
	for i, ap := range id.approx {
		id.approx[i] = constructs.ResolvedTempReference(ap, required)
	}
	for i, ex := range id.exact {
		id.exact[i] = constructs.ResolvedTempReference(ex, required)
	}
}

func (id *interfaceDescImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, id.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, id.Kind(), id.index)
	}
	ab := append(append([]constructs.Abstract{}, id.abstracts...), id.additions...)
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, id.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, id.index).
		AddNonZero(ctx.Short(), `pin`, id.pinnedPkg).
		AddNonZero(ctx.OnlyIndex(), `abstracts`, ab).
		AddNonZero(ctx.Short(), `approx`, id.approx).
		AddNonZero(ctx.Short(), `exact`, id.exact).
		AddNonZero(ctx.OnlyIndex(), `inherits`, id.inherits.ToSlice()).
		AddNonZero(ctx, `hint`, string(id.hint))
}

func (id *interfaceDescImp) ToStringer(s stringer.Stringer) {
	hasExact := len(id.exact) > 0
	hasApprox := len(id.approx) > 0
	hasAbstracts := len(id.abstracts) > 0

	if !hasExact && !hasApprox && !hasAbstracts {
		s.Write(`any`)
		return
	}

	if id.IsPinned() {
		s.Write(id.pinnedPkg.Path(), `:`)
	}

	s.Write(`interface{`)
	next := ``
	if hasExact {
		next = `; `
		s.WriteList(``, `|`, ``, id.exact)
		s.WriteList(`|~`, `|~`, ``, id.approx)
	} else if hasApprox {
		next = `; `
		s.WriteList(`~`, `|~`, ``, id.approx)
	}

	if hasAbstracts {
		s.WriteList(next, `; `, ``, id.abstracts)
	}
	s.Write(` }`)

	check := `interface{ $equal func(other any) bool }`
	if str, ok := strings.CutSuffix(s.String(), check); ok {
		s.Reset().Write(str, `comparable`)
	}
}

func (id *interfaceDescImp) String() string {
	return stringer.String(id)
}
