package constructs

import (
	"go/token"
	"go/types"
	"slices"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/visitor"
)

type (
	Interface interface {
		TypeDesc
		_interface()

		isSupertypeOf(other Interface) bool
		addInheritors(inter Interface) bool
		findImplements(c Class) bool
		setInheritance()
		sortInheritance()
	}

	InterfaceArgs struct {
		RealType   *types.Interface
		Methods    []Named
		TypeParams []Named

		// Exact types are like `string|int|bool` where the
		// data type must match exactly.
		Exact []TypeDesc

		// Approx types are like `~string|~int` where the data type
		// may be exact or an extension of the base type.
		Approx []TypeDesc

		// Package is only needed if the real type is nil
		// so that a Go interface type has to be created.
		Package Package
	}

	interfaceImp struct {
		realType *types.Interface

		typeParams []Named
		methods    []Named
		exact      []TypeDesc
		approx     []TypeDesc

		// FUTURE: Think about tracking all realizations of type parameters
		//         and which exact and approx were actually used in code.
		//         For example, maybe only `map[string]int` ever uses map.

		index      int
		inherits   []Interface
		inheritors []Interface
	}
)

func newInterface(args InterfaceArgs) Interface {
	slices.SortFunc(args.Methods, Compare)
	slices.SortFunc(args.Exact, Compare)
	slices.SortFunc(args.Approx, Compare)

	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(terror.New(`must provide a package if the real type for an interface is nil`))
		}

		mTyp := []*types.Func{}
		pkg := args.Package.Source().Types
		for _, named := range args.Methods {
			name := named.Name()
			sig := named.Type().GoType().(*types.Signature)
			if utils.IsNil(sig) {
				panic(terror.New(`nil signature`).
					With(`name`, name))
			}
			f := types.NewFunc(token.NoPos, pkg, name, sig)
			mTyp = append(mTyp, f)
		}

		tEmb := make([]types.Type, len(args.TypeParams))
		for i, n := range args.TypeParams {
			tEmb[i] = n.GoType()
		}

		if len(args.Exact) > 0 || len(args.Approx) > 0 {
			terms := make([]*types.Term, 0, len(args.Exact)+len(args.Approx))
			for _, ex := range args.Exact {
				terms = append(terms, types.NewTerm(false, ex.GoType()))
			}
			for _, ap := range args.Approx {
				terms = append(terms, types.NewTerm(true, ap.GoType()))
			}
			tEmb = append(tEmb, types.NewUnion(terms))
		}

		realType := types.NewInterfaceType(mTyp, tEmb)
		if realType == nil {
			panic(terror.New(`failed to create an interface`))
		}
		args.RealType = realType
	}

	return &interfaceImp{
		realType:   args.RealType,
		typeParams: args.TypeParams,
		methods:    args.Methods,
		exact:      args.Exact,
		approx:     args.Approx,
	}
}

func (it *interfaceImp) _interface()        {}
func (it *interfaceImp) Kind() kind.Kind    { return kind.Interface }
func (it *interfaceImp) SetIndex(index int) { it.index = index }
func (it *interfaceImp) GoType() types.Type { return it.realType }

func (it *interfaceImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, it.typeParams...)
	visitor.Visit(v, it.methods...)
	visitor.Visit(v, it.exact...)
	visitor.Visit(v, it.approx...)
	visitor.Visit(v, it.inherits...)
}

func (it *interfaceImp) CompareTo(other Construct) int {
	b := other.(*interfaceImp)
	if cmp := CompareSlice(it.exact, b.exact); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(it.approx, b.approx); cmp != 0 {
		return cmp
	}
	if cmp := CompareSlice(it.typeParams, b.typeParams); cmp != 0 {
		return cmp
	}
	return CompareSlice(it.methods, b.methods)
}

func (it *interfaceImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsShort() {
		return jsonify.New(ctx, it.index)
	}

	ctx2 := ctx.HideKind().Short()
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsKindShown(), `kind`, it.Kind()).
		AddIf(ctx, ctx.IsIndexShown(), `index`, it.index).
		AddNonZero(ctx2, `typeParams`, it.typeParams).
		AddNonZero(ctx2, `inherits`, it.inherits).
		AddNonZeroIf(ctx2, ctx.IsInheritorsShown(), `inheritors`, it.inheritors).
		AddNonZero(ctx2, `approx`, it.approx).
		AddNonZero(ctx2, `exact`, it.exact).
		AddNonZero(ctx2, `methods`, it.methods)
}

func (it *interfaceImp) isSupertypeOf(other Interface) bool {
	otherIt, ok := other.GoType().(*types.Interface)
	return ok && types.Implements(it.realType, otherIt)
}

func (it *interfaceImp) addInheritors(other Interface) bool {
	otherImp := other.(*interfaceImp)
	if it == otherImp {
		return true
	}
	if !otherImp.isSupertypeOf(it) {
		return false
	}

	it.inheritors = addInheritors(it.inheritors, other)
	return true
}

func addInheritors(inheritors []Interface, other Interface) []Interface {
	otherImp := other.(*interfaceImp)

	// Tries to home the given other interface into all siblings.
	homed := false
	for _, inheritor := range inheritors {
		if inheritor.addInheritors(otherImp) {
			homed = true
		}
	}
	if homed {
		return inheritors
	}

	// Move super type siblings into the given other interface.
	changed := false
	for i, inheritor := range inheritors {
		if inheritor.isSupertypeOf(otherImp) {
			otherImp.inheritors = append(otherImp.inheritors, inheritor)
			inheritors[i] = nil
			changed = true
		}
	}
	if changed {
		inheritors = utils.RemoveZeros(inheritors)
	}

	// Add the given other interface to this interface.
	return append(inheritors, otherImp)
}

func (it *interfaceImp) findImplements(c Class) bool {
	if !types.Implements(c.GoType(), it.GoType().(*types.Interface)) {
		return false
	}
	homed := false
	for _, inner := range it.inheritors {
		if inner.findImplements(c) {
			homed = true
		}
	}
	if homed {
		return true
	}

	c.addImplement(it)
	return true
}

func (it *interfaceImp) setInheritance() {
	for _, other := range it.inheritors {
		otherInter, ok := other.(*interfaceImp)
		if !ok {
			panic(terror.New(`interfaces in inheritors must be interfaceImps`).
				WithType(`gotten type`, other).
				With(`gotten value`, other))
		}
		otherInter.inherits = append(otherInter.inherits, it)
	}
}

func (it *interfaceImp) sortInheritance() {
	slices.SortFunc(it.inheritors, Compare)
	slices.SortFunc(it.inherits, Compare)
}
