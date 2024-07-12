package constructs

import "github.com/Snow-Gremlin/goToolbox/utils"

type constructSet[T Construct] struct {
	values []T
}

func newTypeSet[T Construct]() *constructSet[T] {
	return &constructSet[T]{}
}

func (s *constructSet[T]) Insert(t T) T {
	for _, t2 := range s.values {
		if t.Equal(t2) {
			return t2
		}
	}
	s.values = append(s.values, t)
	return t
}

func (s *constructSet[T]) SetIndices(index int) int {
	for _, td := range s.values {
		td.SetIndex(index)
		index++
	}
	return index
}

func (s *constructSet[T]) Remove(predict func(Construct) bool) bool {
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
