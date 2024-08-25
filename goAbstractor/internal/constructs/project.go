package constructs

import (
	"github.com/Snow-Gremlin/goToolbox/collections"

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
	MetricsFactory

	// Declarations
	InterfaceDeclFactory
	MethodFactory
	ObjectFactory
	ValueFactory

	// Type Descriptions
	BasicFactory
	InstanceFactory
	InterfaceDescFactory
	SignatureFactory
	StructDescFactory
	TempReferenceFactory
	TypeParamFactory

	Locs() locs.Set
	AllConstructs() collections.Enumerator[Construct]
	FindType(pkgPath, typeName string, panicOnNotFound bool) (Package, TypeDecl, bool)
}
