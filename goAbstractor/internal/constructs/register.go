package constructs

type Register interface {
	RegisterBasic(t Basic) Basic
	RegisterInterface(t Interface) Interface
	RegisterNamed(t Named) Named
	RegisterSignature(t Signature) Signature
	RegisterSolid(t Solid) Solid
	RegisterStruct(t Struct) Struct
	RegisterTypeDefRef(t TypeDefRef) TypeDefRef
	RegisterUnion(t Union) Union
}
