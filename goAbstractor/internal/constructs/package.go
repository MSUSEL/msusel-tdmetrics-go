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
	InitCount() int

	AddImport(p Package) Package
	AddInterfaceDecl(it InterfaceDecl) InterfaceDecl
	AddMethod(m Method) Method
	AddObject(id Object) Object
	AddValue(v Value) Value

	Imports() collections.ReadonlySortedSet[Package]
	InterfaceDecls() collections.ReadonlySortedSet[InterfaceDecl]
	Methods() collections.ReadonlySortedSet[Method]
	Objects() collections.ReadonlySortedSet[Object]
	Values() collections.ReadonlySortedSet[Value]

	Empty() bool
	FindTypeDecl(name string, nest Method) TypeDecl
	FindDecl(name string, nest Method) Declaration
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
