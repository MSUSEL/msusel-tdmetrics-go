package object

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs/kind"
)

type factoryImp struct {
	constructs.FactoryCore[constructs.Object]
}

var _ constructs.Factory = (*factoryImp)(nil)

func New() constructs.ObjectFactory {
	return &factoryImp{FactoryCore: *constructs.NewFactoryCore(kind.Object, Comparer())}
}

func (f *factoryImp) NewObject(args constructs.ObjectArgs) constructs.Object {
	v := f.Add(newObject(args))
	return args.Package.AddObject(v)
}

func (f *factoryImp) Objects() collections.ReadonlySortedSet[constructs.Object] {
	return f.Items().Readonly()
}
