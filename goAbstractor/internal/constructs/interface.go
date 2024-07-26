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

		IsSupertypeOf(other Interface) bool
		AddInheritors(inter Interface) bool
		SetInheritance()
		SortInheritance()
	}

	InterfaceArgs struct {
		RealType   *types.Interface
		Union      Union
		Methods    []Named
		TypeParams []Named

		// Package is only needed if the real type is nil
		// so that a Go interface type has to be created.
		Package Package
	}

	interfaceImp struct {
		realType *types.Interface

		typeParams []Named
		methods    []Named
		union      Union

		index      int
		inherits   []Interface
		inheritors []Interface
	}
)

func newInterface(args InterfaceArgs) Interface {
	methods := slices.Clone(args.Methods)
	tp := slices.Clone(args.TypeParams)

	if utils.IsNil(args.RealType) {
		if utils.IsNil(args.Package) {
			panic(terror.New(`must provide a package if the real type for an interface is nil`))
		}

		mTyp := []*types.Func{}
		pkg := args.Package.Source().Types
		for _, named := range methods {
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
		if !utils.IsNil(args.Union) {
			tEmb = append(tEmb, args.Union.GoType())
		}

		realType := types.NewInterfaceType(mTyp, tEmb)
		if realType == nil {
			panic(terror.New(`failed to create an interface`))
		}
		args.RealType = realType
	}

	return &interfaceImp{
		realType:   args.RealType,
		typeParams: tp,
		methods:    methods,
		union:      args.Union,
	}
}

func (it *interfaceImp) _interface()        {}
func (it *interfaceImp) Kind() kind.Kind    { return kind.Interface }
func (it *interfaceImp) SetIndex(index int) { it.index = index }
func (it *interfaceImp) GoType() types.Type { return it.realType }

func (it *interfaceImp) Visit(v visitor.Visitor) {
	visitor.Visit(v, it.typeParams...)
	visitor.Visit(v, it.methods...)
	visitor.Visit(v, it.union)
	visitor.Visit(v, it.inherits...)
}

func (it *interfaceImp) CompareTo(other Construct) int {
	b := other.(*interfaceImp)
	if cmp := Compare(it.union, b.union); cmp != 0 {
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
		AddNonZero(ctx2, `typeParams`, it.typeParams).
		AddNonZero(ctx2, `inherits`, it.inherits).
		AddNonZero(ctx2, `union`, it.union).
		AddNonZero(ctx2, `methods`, it.methods)
}

func (it *interfaceImp) IsSupertypeOf(other Interface) bool {
	otherIt, ok := other.GoType().(*types.Interface)
	if !ok || utils.IsNil(it.realType) || utils.IsNil(otherIt) {
		// Baked in types don't have underlying interfaces
		// but also shouldn't be needed for any inheritance.
		return false
	}
	return types.Implements(it.realType, otherIt)
}

func (it *interfaceImp) AddInheritors(other Interface) bool {
	otherImp := other.(*interfaceImp)
	if it == otherImp {
		return true
	}
	if !otherImp.IsSupertypeOf(it) {
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
		if inheritor.AddInheritors(otherImp) {
			homed = true
		}
	}
	if homed {
		return inheritors
	}

	// Move super type siblings into the given other interface.
	changed := false
	for i, inheritor := range inheritors {
		if inheritor.IsSupertypeOf(otherImp) {
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

func (it *interfaceImp) SetInheritance() {
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

func (it *interfaceImp) SortInheritance() {
	slices.SortFunc(it.inheritors, Compare)
	slices.SortFunc(it.inherits, Compare)
}
