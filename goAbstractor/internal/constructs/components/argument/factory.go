package argument

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"
)

type ArgumentFactory interface {
	NewArgument(args Args) Argument
	Arguments() collections.ReadonlySet[Argument]
}

type factoryImp struct {
	args collections.SortedSet[Argument]
}

func New() ArgumentFactory {
	return &factoryImp{args: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewArgument(args Args) Argument {
	v, _ := f.args.TryAdd(newArgument(args))
	return v
}

func (f *factoryImp) Arguments() collections.ReadonlySet[Argument] {
	return f.args.Readonly()
}
