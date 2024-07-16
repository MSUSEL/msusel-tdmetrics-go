package locs

import "go/token"

type (
	Set interface {
		NewLoc(p token.Pos) Loc
	}

	setImp struct {
		fs   *token.FileSet
		locs map[token.Pos]Loc
	}
)

func NewSet(fs *token.FileSet) Set {
	return &setImp{
		fs:   fs,
		locs: map[token.Pos]Loc{},
	}
}

func (s *setImp) NewLoc(p token.Pos) Loc {
	if c, ok := s.locs[p]; ok {
		return c
	}
	c := newLoc(s, p)
	s.locs[p] = c
	return c
}
