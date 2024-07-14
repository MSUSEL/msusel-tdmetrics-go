package kind

import "strings"

type Kind string

const (
	Basic     Kind = `basic`
	Class     Kind = `class`
	InterDef  Kind = `interDef`
	Interface Kind = `interface`
	Method    Kind = `method`
	Named     Kind = `named`
	Package   Kind = `package`
	Reference Kind = `reference`
	Signature Kind = `signature`
	Solid     Kind = `solid`
	Struct    Kind = `struct`
	Union     Kind = `union`
	Value     Kind = `value`
)

func (k Kind) CompareTo(other Kind) int {
	return strings.Compare(string(k), string(other))
}
