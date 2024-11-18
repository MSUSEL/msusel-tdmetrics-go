//go:build test

package animals

import "test0014/enums"

type Animal interface {
	Kind() enums.AnimalKind
	isAnimal()
}

func New[B enums.CatBreed | enums.DogBreed](breed B) Animal {
	if !enums.Valid(enums.Enum(breed)) {
		panic(`invalid breed`)
	}
	switch b := any(breed).(type) {
	case enums.CatBreed:
		return cat{breed: b}
	case enums.DogBreed:
		return dog{breed: b}
	}
	panic(`unexpected enum`)
}
