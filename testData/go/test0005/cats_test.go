//go:build test

package cats

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
)

func StrEqual(t *testing.T, val, exp string) {
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

func Test_Cat_Youngest(t *testing.T) {
	c1 := NewCat(`mittens`, 4)
	c2 := NewCat(`missy`, 6)
	c3 := NewCat(`tammy`, 2)
	cats := Cats{c1, c2, c3}
	c4 := cats.Youngest()
	StrEqual(t, c4.String(), `tammy`)
	IntEqual(t, c4.Age, 2)

	NextYear(c3)
	c5 := cats.Youngest()
	StrEqual(t, c5.String(), `tammy`)
	IntEqual(t, c5.Age, 3)

	NextYear(cats...)
	c6 := cats.Youngest()
	StrEqual(t, c6.String(), `tammy`)
	IntEqual(t, c6.Age, 4)

	c2.Age = 2
	c7 := cats.Youngest()
	StrEqual(t, c7.String(), `missy`)
	IntEqual(t, c7.Age, 2)
}

func Test_Cat_Log(t *testing.T) {
	s := &strings.Builder{}
	bckLog := log
	defer func() {
		log = bckLog
	}()
	log = func(a string) {
		fmt.Fprintln(s, a)
	}

	c1 := NewCat(`mittens`, 4)
	c1.Meow()
	Pet(c1)

	StrEqual(t, s.String(), "mittens meows\npetting mittens\n")
}
