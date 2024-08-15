package kind

type Kind string

const (
	Abstract      Kind = `abstract`
	Argument      Kind = `argument`
	Basic         Kind = `basic`
	Field         Kind = `field`
	Instance      Kind = `instance`
	InterfaceDecl Kind = `interfaceDecl`
	InterfaceDesc Kind = `interfaceDesc`
	Method        Kind = `method`
	Object        Kind = `object`
	Package       Kind = `package`
	Reference     Kind = `reference`
	Signature     Kind = `signature`
	StructDesc    Kind = `structDesc`
	TypeParam     Kind = `typeParam`
	Value         Kind = `value`
)
