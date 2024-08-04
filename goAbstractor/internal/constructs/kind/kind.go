package kind

type Kind string

const (
	Basic     Kind = `basic`
	Field     Kind = `field`
	Instance  Kind = `instance`
	Method    Kind = `method`
	Object    Kind = `object`
	Package   Kind = `package`
	Reference Kind = `reference`
	Signature Kind = `signature`
	TypeParam Kind = `typeParam`
	Value     Kind = `value`
)
