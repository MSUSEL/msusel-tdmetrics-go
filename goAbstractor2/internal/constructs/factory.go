package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor2/internal/assert"
)

type ConstructConstraint interface {
	comparable
	Construct
}

type PopulateFunc[Source comparable, T ConstructConstraint] func(*Project, Source, *T)

type Factory[Source comparable, T ConstructConstraint] struct {
	proj       *Project
	associates map[Source]*T
	instances  map[T]*T
	populate   PopulateFunc[Source, T]
}

func NewFactory[Source comparable, T ConstructConstraint](proj *Project, populate PopulateFunc[Source, T]) *Factory[Source, T] {
	assert.NotNil(proj, `project for function factory`)
	assert.NotNil(populate, `populate func for function factory`)
	return &Factory[Source, T]{
		proj:       proj,
		associates: map[Source]*T{},
		instances:  map[T]*T{},
		populate:   populate,
	}
}

func (f Factory[Source, T]) Kind() string { return utils.Zero[T]().Kind() }

func (f *Factory[Source, T]) New(src Source) *T {
	assert.NotNil(f.associates, `factory must be created with NewFactory`)
	assert.NotNil(src, `source in factory new`)

	if v, found := f.associates[src]; found {
		return v
	}

	v := new(T)
	f.associates[src] = v

	func() {
		defer func() {
			if r := recover(); r != nil {
				delete(f.associates, src)
				panic(terror.RecoveredPanic(r))
			}
		}()
		f.populate(f.proj, src, v)
	}()

	if other, found := f.instances[*v]; found {
		f.associates[src] = other
		return other
	}

	f.instances[*v] = v
	return v
}
