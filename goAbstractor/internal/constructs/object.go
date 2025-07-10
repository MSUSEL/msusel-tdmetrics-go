package constructs

import (
	"go/types"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/locs"
	"github.com/Snow-Gremlin/goToolbox/collections"
)

// Object is a named type typically explicitly defined at the given location
// in the source code. An object typically handles structs with optional
// parameter types. An object can handle any type that methods can use
// as a receiver.
//
// If type parameters are given then the object is generic.
// Instances with realized versions of the object,
// are added for each used instance in the source code.
// If there are no instances then the generic object isn't used.
// The type parameters will also contain any nesting type parameters
// from any generic methods enclosing this object.
type Object interface {
	TypeDecl
	Nestable
	IsObject()

	Data() StructDesc
	Methods() collections.ReadonlySortedSet[Method]
	Interface() InterfaceDesc

	AddMethod(met Method) Method
	SetInterface(it InterfaceDesc)

	IsNamed() bool
	IsGeneric() bool
	IsNested() bool
	TypeParams() []TypeParam
	ImplicitTypeParams() []TypeParam
	AddInstance(inst ObjectInst) ObjectInst
	Instances() collections.ReadonlySortedSet[ObjectInst]
	FindInstance(implicitTypes, instanceTypes []TypeDesc) (TypeDesc, bool)
}

type ObjectArgs struct {
	RealType types.Type
	Package  Package
	Name     string
	Exported bool
	Location locs.Loc
	Nest     NestType

	TypeParams []TypeParam
	Data       StructDesc
}

type ObjectFactory interface {
	NewObject(args ObjectArgs) Object
	Objects() collections.ReadonlySortedSet[Object]
}
