package constructs

import "github.com/Snow-Gremlin/goToolbox/collections"

// MethodInst represents an instantiation of generic method
// has been resolved to a specific method with specific type parameters,
// e.g. List[T any] might be resolved to List<int>.
//
// The instance parameter may be referencing a type parameter,
// e.g. List[T any] might be resolved to List<S int|string>, thus the instance
// is also generic. The type may be a non-type parameter on a generic type,
// e.g. List[List[T any]] where List[T any] is the instance type.
type MethodInst interface {
	NestType
	Construct
	TempReferenceContainer
	IsMethodInst()

	Generic() Method
	Resolved() Signature
	InstanceTypes() []TypeDesc

	Receiver() ObjectInst
	SetReceiver(obj ObjectInst)
}

type MethodInstArgs struct {
	Generic       Method
	Resolved      Signature
	InstanceTypes []TypeDesc
}

type MethodInstFactory interface {
	NewMethodInst(args MethodInstArgs) MethodInst
	MethodInsts() collections.ReadonlySortedSet[MethodInst]
}
