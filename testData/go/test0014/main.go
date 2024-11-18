//go:build test

package main

import (
	"test0014/animals"
	"test0014/enums"
)

// A test for when there are unexported abstracts on interfaces causing
// the interface to be locked to a package. Requires multiple packages
// to test the locking so we'll check that too and module file.

var pets = [3]animals.Animal{
	animals.New(enums.Husky),
	animals.New(enums.Poodle),
	animals.New(enums.MaineCoon),
}

func main() {
	for _, pet := range pets {
		switch pet.Kind() {
		case enums.Cat:
			println(any(pet).(animals.Cat).Breed())
		case enums.Dog:
			println(any(pet).(animals.Dog).Breed())
		default:
			println(`unknown animal`)
		}
	}
}
