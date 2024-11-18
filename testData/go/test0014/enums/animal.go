//go:build test

package enums

type AnimalKind string

const (
	Cat AnimalKind = `Cat`
	Dog AnimalKind = `Dog`
)

func (a AnimalKind) valid() bool {
	switch a {
	case Cat, Dog:
		return true
	}
	return false
}
