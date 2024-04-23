//go:build test

package cats

import (
	"strconv"
	"strings"
	"testing"
)

func StrEqual(t *testing.T, val string, expLines ...string) {
	exp := strings.Join(expLines, "\n")
	if val != exp {
		t.Error("Unexpected string:" +
			"\n   Value:    " + val +
			"\n   Expected: " + exp)
	}
}

func IntEqual(t *testing.T, val, exp int) {
	if val != exp {
		t.Error("Unexpected string:" +
			"\n   Value:    " + strconv.Itoa(val) +
			"\n   Expected: " + strconv.Itoa(exp))
	}
}

func Test_CatTable(t *testing.T) {
	ct := CatTable{}
	ct.AddNew(`mittens`, Persian, 4)
	ct.AddNew(`missy`, Himalayan, 6)
	ct.AddNew(`tammy`, Siberian, 2)
	ct.AddNew(`brat`, Himalayan, 3)
	ct.AddNew(`missy`, Bengal, 4)

	StrEqual(t, ct.String(),
		`|ID   |Name   |Breed    |Age|`,
		`|:----|:------|:--------|--:|`,
		`|cat-0|mittens|Persian  |  4|`,
		`|cat-1|missy  |Himalayan|  6|`,
		`|cat-2|tammy  |Siberian |  2|`,
		`|cat-3|brat   |Himalayan|  3|`,
		`|cat-4|missy  |Bengal   |  4|`)

	StrEqual(t, ct.CatsWithName(`mittens`).String(), `cat-0`)
	StrEqual(t, ct.CatsWithName(`missy`).String(), `cat-1, cat-4`)

	StrEqual(t, ct.CatsWithBreed(Persian).String(), `cat-0`)
	StrEqual(t, ct.CatsWithBreed(Himalayan).String(), `cat-1, cat-3`)

	StrEqual(t, ct.AllIDs().String(), `cat-0, cat-1, cat-2, cat-3, cat-4`)

	id := ID(`cat-2`)
	StrEqual(t, ct.Name(id), `tammy`)
	StrEqual(t, ct.Breed(id).String(), `Siberian`)
	IntEqual(t, ct.Age(id), 2)

	min, max := ct.AgeRange()
	IntEqual(t, min, 2)
	IntEqual(t, max, 6)

	breeds := ct.AllBreeds()
	strBreeds := make([]string, len(breeds))
	for i, breed := range breeds {
		strBreeds[i] = breed.String()
	}
	StrEqual(t, strings.Join(strBreeds, ", "), `Persian, Bengal, Siberian, Himalayan`)
}
