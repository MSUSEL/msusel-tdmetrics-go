//go:build test

package animals

import "test0014/enums"

type Dog interface {
	Breed() enums.DogBreed
	isDog()
}

type dog struct {
	breed enums.DogBreed
}

func (d dog) Kind() enums.AnimalKind { return enums.Dog }
func (d dog) Breed() enums.DogBreed  { return d.breed }

func (d dog) isDog()    {}
func (d dog) isAnimal() {}
