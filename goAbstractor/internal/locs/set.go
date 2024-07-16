package locs

import (
	"go/token"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Set interface {
		NewLoc(p token.Pos) Loc
		Reset()
		flag(p token.Pos)
		indexFor(p token.Pos) int
		infoFor(p token.Pos) (int, string, int)
	}

	setImp struct {
		fs       *token.FileSet
		locs     map[token.Pos]Loc
		flagged  map[token.Pos]bool
		indices  map[token.Pos]int
		files    []string
		lines    []int
		finished bool
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

func (s *setImp) Reset() {
	s.flagged = map[token.Pos]bool{}
	s.indices = map[token.Pos]int{}
	s.files = []string{}
	s.lines = []int{}
	s.finished = false
}

func (s *setImp) flag(p token.Pos) {
	if s.finished {
		panic(terror.New(`flagging a location must be after a reset ` +
			`and prior to any location information looked up`))
	}
	s.flagged[p] = true
}

func (s *setImp) finish() {
	if s.finished {
		return
	}

	files := map[string]int{}
	for p := range s.flagged {
		f := s.fs.File(p)
		files[f.Name()] = f.LineCount()
	}

	// TODO: Implement

	s.finished = true
}

func (s *setImp) indexFor(p token.Pos) int {
	s.finish()
	return s.indices[p]
}

func (s *setImp) infoFor(p token.Pos) (int, string, int) {
	s.finish()
	index := s.indices[p]
	name, line := ``, 0
	if index > 0 {
		fsp := s.fs.Position(p)
		name, line = fsp.Filename, fsp.Line
	}
	return index, name, line
}

func (s *setImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	s.finish()
	m := jsonify.NewMap()
	line := 1
	for i, file := range s.files {
		m.Add(ctx, strconv.Itoa(line), file)
		line += s.lines[i]
	}
	return m
}
