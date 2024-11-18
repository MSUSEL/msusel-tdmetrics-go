//go:build test

package enums

type Enum interface {
	valid() bool
}

func Valid(e Enum) bool {
	return e.valid()
}
