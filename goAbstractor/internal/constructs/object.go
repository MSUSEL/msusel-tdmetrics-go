package constructs

// Object is a named type typically explicitly defined at the given location
// in the source code. An object typically handles structs with optional
// parameter types. An object can handle any type that methods can use
// as a receiver.
//
// If type parameters are given then the object is generic.
// Instances with realized versions of the object,
// are added for each used instance in the source code.
// If there are no instances then the generic object isn't used.
type Object interface {
	Declaration
	TypeDesc
	IsObject()

	AddMethod(met Method) Method
	AddInstance(inst Instance) Instance
	SetInterface(it InterfaceDesc)

	IsNamed() bool
	IsGeneric() bool
}
