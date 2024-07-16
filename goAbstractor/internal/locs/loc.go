package locs

import "go/token"

type (
	Loc interface {
		_loc()
	}

	locImp struct {
		s Set
		p token.Pos
	}
)

func newLoc(s Set, p token.Pos) Loc {
	return &locImp{s: s, p: p}
}

func (c *locImp) _loc() {}
