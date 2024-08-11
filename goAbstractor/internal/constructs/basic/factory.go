package basic

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	basics collections.SortedSet[constructs.Basic]
}

func New() constructs.BasicFactory {
	return &factoryImp{basics: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewBasic(args constructs.BasicArgs) constructs.Basic {
	v, _ := f.basics.TryAdd(newBasic(args))
	return v
}

func (f *factoryImp) Basics() collections.ReadonlySet[constructs.Basic] {
	return f.basics.Readonly()
}
