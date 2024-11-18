//go:build test

package enums

type DogBreed string

const (
	Poodle    DogBreed = `Poodle`
	Chihuahua DogBreed = `Chihuahua`
	Husky     DogBreed = `Husky`
)

func (c DogBreed) valid() bool {
	switch c {
	case Poodle, Chihuahua, Husky:
		return true
	}
	return false
}
