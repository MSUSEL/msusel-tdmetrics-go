//go:build test

package animals

import "test0014/enums"

type Cat interface {
	Breed() enums.CatBreed
	isCat()
}

type cat struct {
	breed enums.CatBreed
}

func (c cat) Kind() enums.AnimalKind { return enums.Cat }
func (c cat) Breed() enums.CatBreed  { return c.breed }

func (c cat) isCat()    {}
func (c cat) isAnimal() {}
