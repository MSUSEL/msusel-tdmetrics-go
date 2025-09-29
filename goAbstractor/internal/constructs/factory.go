package constructs

import (
	"fmt"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/enumerator"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
	"github.com/Snow-Gremlin/goToolbox/comp"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type Factory interface {
	Kind() kind.Kind
	Enumerate() collections.Enumerator[Construct]
	Refresh()
	String() string
}

type FactoryCore[T Construct] struct {
	kind  kind.Kind
	items collections.SortedSet[T]
}

var _ Factory = (*FactoryCore[Abstract])(nil)

func NewFactoryCore[T Construct](kind kind.Kind, comparer comp.Comparer[T]) *FactoryCore[T] {
	return &FactoryCore[T]{
		kind:  kind,
		items: sortedSet.New(comparer),
	}
}

func (f *FactoryCore[T]) Kind() kind.Kind { return kind.Abstract }

func (f *FactoryCore[T]) Add(item T) T {
	v, _ := f.items.TryAdd(item)
	return v
}

func (f *FactoryCore[T]) Items() collections.SortedSet[T] { return f.items }

func (f *FactoryCore[T]) Enumerate() collections.Enumerator[Construct] {
	return enumerator.Cast[Construct](f.items.Enumerate())
}

func (f *FactoryCore[T]) Refresh() { f.items.Refresh() }

func (f *FactoryCore[T]) String() string {
	buf := &strings.Builder{}
	buf.WriteString(f.Kind().Plural())
	if f.Enumerate().Empty() {
		buf.WriteString(" { }")
		return buf.String()
	}
	buf.WriteString(" {\n")
	i := 0
	for c := range f.Enumerate().Seq() {
		extra, state := ``, ``
		if !c.Alive() {
			state += `X`
		}
		if c.Duplicate() {
			state += `D`
		}
		extra = fmt.Sprintf(`[%s%2d]`, state, c.Index())
		fmt.Fprintf(buf, "  %2d. %s%q\n", i+1, extra, c.String())
		i++
	}
	buf.WriteString("}")
	return buf.String()
}
