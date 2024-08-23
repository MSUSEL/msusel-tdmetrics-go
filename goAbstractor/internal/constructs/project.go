package constructs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
)

type Project interface {
	jsonify.Jsonable

	// Components
	AbstractFactory
	ArgumentFactory
	FieldFactory
	PackageFactory

	// Declarations
	InterfaceDeclFactory
	MethodFactory
	ObjectFactory
	ValueFactory

	// Type Descriptions
	BasicFactory
	InstanceFactory
	InterfaceDescFactory
	ReferenceFactory
	SignatureFactory
	StructDescFactory
	TypeParamFactory

	NewLoc(pos token.Pos) locs.Loc
	FindType(pkgPath, typeName string, panicOnNotFound bool) (Package, TypeDecl, bool)
}
