package constructs

import "github.com/Snow-Gremlin/goToolbox/utils"

type typeSet[T TypeDesc] struct {
	values []T
}

func newTypeSet[T TypeDesc]() *typeSet[T] {
	return &typeSet[T]{}
}

func (s *typeSet[T]) Insert(t T) T {
	for _, t2 := range s.values {
		if t.Equal(t2) {
			return t2
		}
	}
	s.values = append(s.values, t)
	return t
}

func (s *typeSet[T]) SetIndices(index int) int {
	for _, td := range s.values {
		td.SetIndex(index)
		index++
	}
	return index
}

func (s *typeSet[T]) Remove(predict func(TypeDesc) bool) bool {
	changed := false
	zero := utils.Zero[T]()
	for i, td := range s.values {
		if predict(td) {
			changed = true
			s.values[i] = zero
		}
	}
	s.values = utils.RemoveZeros(s.values)
	return changed
}
