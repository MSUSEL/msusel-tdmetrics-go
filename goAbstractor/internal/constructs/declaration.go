package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"

type Declaration interface {
	Construct
	IsDeclaration()

	Package() Package
	Name() string
	Location() locs.Loc
}

type TypeDecl interface {
	Declaration
	TypeDesc
}
