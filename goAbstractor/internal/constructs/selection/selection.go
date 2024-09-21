package selection

import (
	"strings"

	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/assert"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type selectionImp struct {
	name   string
	origin constructs.Construct
	index  int
	alive  bool
}

func newSelection(args constructs.SelectionArgs) constructs.Selection {
	assert.ArgValidId(`name`, args.Name)
	assert.ArgNotNil(`origin`, args.Origin)
	return &selectionImp{
		name:   args.Name,
		origin: args.Origin,
	}
}

func (s *selectionImp) IsSelection() {}

func (s *selectionImp) Kind() kind.Kind     { return kind.Selection }
func (s *selectionImp) Index() int          { return s.index }
func (s *selectionImp) SetIndex(index int)  { s.index = index }
func (s *selectionImp) Alive() bool         { return s.alive }
func (s *selectionImp) SetAlive(alive bool) { s.alive = alive }
func (s *selectionImp) Name() string        { return s.name }

func (s *selectionImp) Origin() constructs.Construct { return s.origin }

func (s *selectionImp) RemoveTempReferences() {
	if s.origin.Kind() == kind.TempReference {
		s.origin = constructs.ResolvedTempReference(s.origin.(constructs.TempReference))
	}
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
		return jsonify.New(ctx, s.index)
	}
	if ctx.IsShort() {
		return jsonify.NewSprintf(`%s%d`, s.Kind(), s.index)
	}
	return jsonify.NewMap().
		AddIf(ctx, ctx.IsDebugKindIncluded(), `kind`, s.Kind()).
		AddIf(ctx, ctx.IsDebugIndexIncluded(), `index`, s.index).
		Add(ctx, `name`, s.name).
		Add(ctx.Short(), `origin`, s.origin)
}

func (s *selectionImp) String() string {
	buf := &strings.Builder{}
	buf.WriteString(s.origin.String())
	buf.WriteString(`.`)
	buf.WriteString(s.name)
	return buf.String()
}
