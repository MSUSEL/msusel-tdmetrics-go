package interfaceDecl

import (
	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/sortedSet"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/constructs"
)

type factoryImp struct {
	methods collections.SortedSet[constructs.InterfaceDecl]
}

func New() constructs.InterfaceDeclFactory {
	return &factoryImp{methods: sortedSet.New(Comparer())}
}

func (f *factoryImp) NewInterfaceDecl(args constructs.InterfaceDeclArgs) constructs.InterfaceDecl {
	v, _ := f.methods.TryAdd(newInterfaceDecl(args))
	return args.Package.AddInterfaceDecl(v)
}

func (f *factoryImp) InterfaceDecls() collections.ReadonlySortedSet[constructs.InterfaceDecl] {
	return f.methods.Readonly()
}
