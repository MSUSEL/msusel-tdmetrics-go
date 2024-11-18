//go:build test

package enums

type CatBreed string

const (
	Siamese   CatBreed = `Siamese`
	MaineCoon CatBreed = `MaineCoon`
	Persian   CatBreed = `Persian`
)

func (c CatBreed) valid() bool {
	switch c {
	case Siamese, MaineCoon, Persian:
		return true
	}
	return false
}
