package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

// Declaration is a type, value, or method declaration with a name.
type Declaration interface {
	Construct

	// IsDeclaration indicates that the type is a Declaration at compile time.
	// This prevents anything else from duck-typing into a Declaration.
	IsDeclaration()

	TypeParams() []TypeParam
	AddInstance(inst Instance) Instance

	Package() Package
	Name() string
	Location() locs.Loc
}

// TypeDecl is both a type description and a type declaration,
// i.e. `type Foo struct{}; var X Foo`.
type TypeDecl interface {
	Declaration
	TypeDesc
}
