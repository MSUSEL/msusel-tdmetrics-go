package kind

type Kind string

const (
	Basic       Kind = `basic`
	Declaration Kind = `declaration`
	Field       Kind = `field`
	Instance    Kind = `instance`
	Method      Kind = `method`

	Package   Kind = `package`
	Reference Kind = `reference`
	Signature Kind = `signature`
	TypeParam Kind = `typeParam`
	Union     Kind = `union`
	Value     Kind = `value`
)
