package kind

type Kind string

const (
	Abstract         Kind = `abstract`
	Argument         Kind = `argument`
	Basic            Kind = `basic`
	Field            Kind = `field`
	InterfaceDecl    Kind = `interfaceDecl`
	InterfaceDesc    Kind = `interfaceDesc`
	InterfaceInst    Kind = `interfaceInst`
	Method           Kind = `method`
	MethodInst       Kind = `methodInst`
	Metrics          Kind = `metrics`
	Object           Kind = `object`
	ObjectInst       Kind = `objectInst`
	Package          Kind = `package`
	Selection        Kind = `selection`
	Signature        Kind = `signature`
	StructDesc       Kind = `structDesc`
	TempDeclRef      Kind = `tempDeclRef`
	TempReference    Kind = `tempReference`
	TempTypeParamRef Kind = `tempTypeParamRef`
	TypeParam        Kind = `typeParam`
	Value            Kind = `value`
)
