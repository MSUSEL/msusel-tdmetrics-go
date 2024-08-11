package interfaceDecl

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	methods collections.SortedSet[constructs.InterfaceDecl]
}

func NewFactory() constructs.InterfaceDeclFactory {
	return &factoryImp{methods: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInterfaceDecl(args constructs.InterfaceDeclArgs) constructs.InterfaceDecl {
	v, _ := f.methods.TryAdd(newInterfaceDecl(args))
	return v
}

func (f *factoryImp) InterfaceDecls() collections.ReadonlySet[constructs.InterfaceDecl] {
	return f.methods.Readonly()
}
