package constructs

import (
	"golang.org/x/tools/go/packages"

	"github.com/Snow-Gremlin/goToolbox/collections"
)

type Package interface {
	Construct
	IsPackage()

	Source() *packages.Package
	Path() string
	Name() string
	EntryPoint() bool
	ImportPaths() []string
	Imports() collections.ReadonlySortedSet[Package]
	InitCount() int

	AddImport(p Package) Package
	AddInterfaceDecl(it InterfaceDecl) InterfaceDecl
	AddMethod(m Method) Method
	AddObject(id Object) Object
	AddValue(v Value) Value

	Empty() bool
	FindTypeDecl(name string) TypeDecl
	ResolveReceivers()
}

type PackageArgs struct {
	RealPkg     *packages.Package
	Path        string
	Name        string
	ImportPaths []string
}

type PackageFactory interface {
	NewPackage(args PackageArgs) Package
	Packages() collections.ReadonlySortedSet[Package]
	FindPackageByPath(path string) Package
}
