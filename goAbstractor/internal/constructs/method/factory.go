package method

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	methods collections.SortedSet[constructs.Method]
}

func New() constructs.MethodFactory {
	return &factoryImp{methods: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewMethod(args constructs.MethodArgs) constructs.Method {
	v, _ := f.methods.TryAdd(newMethod(args))
	return args.Package.AddMethod(v)
}

func (f *factoryImp) Methods() collections.ReadonlySortedSet[constructs.Method] {
	return f.methods.Readonly()
}
