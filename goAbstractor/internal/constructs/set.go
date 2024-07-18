package constructs

import (
	"sort"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/list"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	UntypedSet interface {
		SetIndices(index int) int
		Remove(predicate func(Construct) bool) bool
	}

	Set[T Construct] interface {
		UntypedSet
		Values() collections.ReadonlyList[T]
		Insert(t T) T
	}

	setImp[T Construct] struct {
		values collections.List[T]
	}
)

func NewSet[T Construct]() Set[T] {
	return &setImp[T]{
		values: list.New[T](),
	}
}

func (s *setImp[T]) Values() collections.ReadonlyList[T] {
	return s.values.Readonly()
}

func (s *setImp[T]) Insert(t T) T {
	if utils.IsNil(t) {
		panic(terror.New(`may not insert a nil value`))
	}

	index, found := sort.Find(s.values.Count(), func(i int) int {
		return t.CompareTo(s.values.Get(i))
	})
	if found {
		return s.values.Get(index)
	}

	s.values.Insert(index, t)
	return t
}

func (s *setImp[T]) SetIndices(index int) int {
	for i := range s.values.Count() {
		s.values.Get(i).SetIndex(index)
		index++
	}
	return index
}

func (s *setImp[T]) Remove(p func(Construct) bool) bool {
	removed := s.values.RemoveIf(func(value T) bool {
		return p(value)
	})
	return removed
}

func (s *setImp[T]) ToJson(ctx *jsonify.Context) jsonify.Datum {
	list := jsonify.NewList()
	s.values.Enumerate().Foreach(func(value T) {
		list.Append(ctx, value)
	})
	return list
}
