package abstract

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	args collections.SortedSet[constructs.Abstract]
}

func New() constructs.AbstractFactory {
	return &factoryImp{args: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewAbstract(args constructs.AbstractArgs) constructs.Abstract {
	v, _ := f.args.TryAdd(newAbstract(args))
	return v
}

func (f *factoryImp) Abstracts() collections.ReadonlySortedSet[constructs.Abstract] {
	return f.args.Readonly()
}
