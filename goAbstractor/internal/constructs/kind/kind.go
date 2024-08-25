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
	Metrics       Kind = `metrics`
	Object        Kind = `object`
	Package       Kind = `package`
	Signature     Kind = `signature`
	StructDesc    Kind = `structDesc`
	TempReference Kind = `tempReference`
	TypeParam     Kind = `typeParam`
	Value         Kind = `value`
)
