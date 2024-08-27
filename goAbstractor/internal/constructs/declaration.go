package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

// Declaration is a type, value, or method declaration with a name.
type Declaration interface {
	Construct

	// IsDeclaration indicates that the type is a Declaration at compile time.
	// This prevents anything else from duck-typing into a Declaration.
	IsDeclaration()

	Package() Package
	Name() string
	Location() locs.Loc
	Type() TypeDesc
}

// TypeDecl is both a type description and a type declaration,
// i.e. `type Foo struct{}; var X Foo`.
type TypeDecl interface {
	Declaration
	TypeDesc
}

/*
func FindInstance(decl Declaration, instanceTypes []TypeDesc) (Instance, bool) {
	cmp := SliceComparer[TypeDesc]()
	return decl.Instances().Enumerate().Where(func(i Instance) bool {
		return cmp(instanceTypes, i.InstanceTypes()) == 0
	}).First()
}
*/
