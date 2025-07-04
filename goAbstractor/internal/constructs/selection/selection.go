package selection

import (
	"github.com/Snow-Gremlin/goToolbox/comp"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/stringer"
)

type selectionImp struct {
	constructs.ConstructCore
	name   string
	origin constructs.Construct
	target constructs.Construct
}

func newSelection(args constructs.SelectionArgs) constructs.Selection {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`origin`, args.Origin)
	return &selectionImp{
		name:   args.Name,
		origin: args.Origin,
		target: args.Target,
	}
}

func (s *selectionImp) IsSelection() {}

func (s *selectionImp) Kind() kind.Kind { return kind.Selection }
func (s *selectionImp) Name() string    { return s.name }

func (s *selectionImp) Origin() constructs.Construct { return s.origin }

func (s *selectionImp) Target() constructs.Construct {
	if utils.IsNil(s.target) {
		return s.target
	}

	switch s.origin.Kind() {
	case kind.InterfaceDecl:
		s.target = setTargetFromInterfaceDecl(s.origin.(constructs.InterfaceDecl), s.name)
	case kind.InterfaceDesc:
		s.target = setTargetFromInterfaceDesc(s.origin.(constructs.InterfaceDesc), s.name)
	case kind.InterfaceInst:
		s.target = setTargetFromInterfaceInst(s.origin.(constructs.InterfaceInst), s.name)
	case kind.Object:
		s.target = setTargetFromObject(s.origin.(constructs.Object), s.name)
	case kind.ObjectInst:
		s.target = setTargetFromObjectInst(s.origin.(constructs.ObjectInst), s.name)
	case kind.Package:
		s.target = setTargetFromPackage(s.origin.(constructs.Package), s.name)
	case kind.StructDesc:
		s.target = setTargetFromStructDesc(s.origin.(constructs.StructDesc), s.name)
	}
	return s.target
}

func setTargetFromInterfaceDecl(t constructs.InterfaceDecl, name string) constructs.Construct {
	return setTargetFromInterfaceDesc(t.Interface(), name)
}

func setTargetFromInterfaceDesc(t constructs.InterfaceDesc, name string) constructs.Construct {
	for _, a := range t.Abstracts() {
		if a.Name() == name {
			return a
		}
	}
	return nil
}

func setTargetFromInterfaceInst(t constructs.InterfaceInst, name string) constructs.Construct {
	return setTargetFromInterfaceDesc(t.Resolved(), name)
}

func setTargetFromObject(t constructs.Object, name string) constructs.Construct {
	for i := range t.Methods().Count() {
		if m := t.Methods().Get(i); m.Name() == name {
			return m
		}
	}
	return setTargetFromStructDesc(t.Data(), name)
}

func setTargetFromObjectInst(t constructs.ObjectInst, name string) constructs.Construct {
	for i := range t.Methods().Count() {
		if m := t.Methods().Get(i); m.Name() == name {
			return m
		}
	}
	return setTargetFromStructDesc(t.ResolvedData(), name)
}

func setTargetFromPackage(t constructs.Package, name string) constructs.Construct {
	for i := range t.Methods().Count() {
		if m := t.Methods().Get(i); !m.HasReceiver() && m.Name() == name {
			return m
		}
	}
	for i := range t.Objects().Count() {
		if obj := t.Objects().Get(i); obj.Name() == name {
			return obj
		}
	}
	for i := range t.InterfaceDecls().Count() {
		if it := t.InterfaceDecls().Get(i); it.Name() == name {
			return it
		}
	}
	for i := range t.Values().Count() {
		if v := t.Values().Get(i); v.Name() == name {
			return v
		}
	}
	return nil
}

func setTargetFromStructDesc(t constructs.StructDesc, name string) constructs.Construct {
	for _, f := range t.Fields() {
		if f.Name() == name {
			return f
		}
	}
	return nil
}

func (s *selectionImp) RemoveTempReferences(required bool) bool {
	changed := false
	if s.origin.Kind() == kind.TempReference {
		s.origin, changed = constructs.ResolvedTempReference(s.origin.(constructs.TempReference), required)
	}
	return changed
}

func (s *selectionImp) CompareTo(other constructs.Construct) int {
	return constructs.CompareTo[constructs.Selection](s, other, Comparer())
}

func Comparer() comp.Comparer[constructs.Selection] {
	return func(a, b constructs.Selection) int {
		aImp, bImp := a.(*selectionImp), b.(*selectionImp)
		if aImp == bImp {
			return 0
		}
		return comp.Or(
			comp.DefaultPend(aImp.name, bImp.name),
			constructs.ComparerPend(aImp.origin, bImp.origin),
		)
	}
}

func (s *selectionImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsOnlyIndex() {
		return jsonify.New(ctx, s.Index())
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, s.Kind(), s.Index())
	}
	if ctx.SkipDead() && !s.Alive() {
		return nil
	}
	if !ctx.KeepDuplicates() && s.Duplicate() {
		return nil
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, s.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, s.Index()).
		AddIf(ctx, ctx.IsDebugAliveIncluded(), `alive`, s.Alive()).
		Add(ctx, `name`, s.name).
		Add(ctx.Short(), `origin`, s.origin).
		AddNonZero(ctx.Short(), `target`, s.Target())
}

func (s *selectionImp) ToStringer(str stringer.Stringer) {
	str.Write(s.origin, `.`, s.name)
	if target := s.Target(); !utils.IsNil(target) {
		str.Write(`=>`, target)
	}
}

func (s *selectionImp) String() string {
	return stringer.String(s)
}
