package kind

type Kind string

const (
	Argument      Kind = `argument`
	Basic         Kind = `basic`
	Field         Kind = `field`
	Instance      Kind = `instance`
	InterfaceDecl Kind = `interfaceDecl`
	InterfaceDesc Kind = `interfaceDesc`
	Method        Kind = `method`
	Object        Kind = `object`
	PackageCon    Kind = `packageCon`
	Reference     Kind = `reference`
	Signature     Kind = `signature`
	StructDesc    Kind = `structDesc`
	TypeDesc      Kind = `typeDesc`
	TypeParam     Kind = `typeParam`
	Value         Kind = `value`
)
