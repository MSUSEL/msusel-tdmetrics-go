package kind

type Kind string

const (
	Abstract      Kind = `abstract`
	Argument      Kind = `argument`
	Basic         Kind = `basic`
	Field         Kind = `field`
	InterfaceDecl Kind = `interfaceDecl`
	InterfaceDesc Kind = `interfaceDesc`
	InterfaceInst Kind = `interfaceInst`
	Method        Kind = `method`
	MethodInst    Kind = `methodInst`
	Metrics       Kind = `metrics`
	Object        Kind = `object`
	ObjectInst    Kind = `objectInst`
	Package       Kind = `package`
	Signature     Kind = `signature`
	StructDesc    Kind = `structDesc`
	TempReference Kind = `tempReference`
	TypeParam     Kind = `typeParam`
	Usage         Kind = `usage`
	Value         Kind = `value`
)
