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

	// UpdateIndices should be called after all types have been registered
	// and all packages have been processed. This will update all the index
	// fields that will be used as references in the output models.
	UpdateIndices()

	ResolveImports()
	ResolveReceivers()
	ResolveInheritance()
	ResolveReferences()
	FlagLocations()
}
