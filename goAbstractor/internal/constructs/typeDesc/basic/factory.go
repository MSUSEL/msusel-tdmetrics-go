package basic

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type BasicFactory interface {
	NewBasic(args Args) Basic
	Basics() collections.ReadonlySet[Basic]
}

type factoryImp struct {
	basics collections.SortedSet[Basic]
}

func New() BasicFactory {
	return &factoryImp{basics: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewBasic(args Args) Basic {
	v, _ := f.basics.TryAdd(newBasic(args))
	return v
}

func (f *factoryImp) Basics() collections.ReadonlySet[Basic] {
	return f.basics.Readonly()
}
