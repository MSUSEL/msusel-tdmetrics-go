package kind

import "strings"

type Kind string

const (
	Basic         Kind = `basic`
	ClassDecl     Kind = `classDecl`
	InterfaceDecl Kind = `interfaceDecl`
	Interface     Kind = `interface`
	Method        Kind = `method`
	Named         Kind = `named`
	Package       Kind = `package`
	Reference     Kind = `reference`
	Signature     Kind = `signature`
	Solid         Kind = `solid`
	Struct        Kind = `struct`
	Union         Kind = `union`
	ValueDecl     Kind = `valueDecl`
)

func (k Kind) CompareTo(other Kind) int {
	return strings.Compare(string(k), string(other))
}
