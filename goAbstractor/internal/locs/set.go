package locs

import (
	"go/token"
	"path/filepath"
	"strconv"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Set interface {
		NewLoc(p token.Pos) Loc
		Reset()
		flag(p token.Pos)
		infoFor(p token.Pos) (int, string, int)
	}

	setImp struct {
		fs       *token.FileSet
		basePath string
		flagged  map[token.Pos]bool
		offsets  map[string]int
		finished bool
	}
)

func NewSet(fs *token.FileSet, basePath string) Set {
	s := &setImp{
		fs:       fs,
		basePath: basePath,
	}
	s.Reset()
	return s
}

func (s *setImp) NewLoc(p token.Pos) Loc {
	return newLoc(s, p)
}

func (s *setImp) Reset() {
	s.flagged = map[token.Pos]bool{}
	s.finished = false
}

func (s *setImp) flag(p token.Pos) {
	if s.finished {
		panic(terror.New(`flagging a location must be after a reset ` +
			`and prior to any location information looked up`))
	}
	s.flagged[p] = true
}

func (s *setImp) relPath(file string) string {
	target, err := filepath.Rel(s.basePath, file)
	if err != nil {
		panic(terror.New(`error creating a relative path for a location`, err).
			With(`base path`, s.basePath).
			With(`file`, file))
	}
	return target
}

func (s *setImp) finish() {
	if s.finished {
		return
	}
	s.finished = true

	lineCounts := map[string]int{}
	for p := range s.flagged {
		f := s.fs.File(p)
		file, lines := s.relPath(f.Name()), f.LineCount()
		lineCounts[file] = lines
	}
	files := utils.SortedKeys(lineCounts)

	s.offsets = map[string]int{}
	offset := 1
	for _, file := range files {
		s.offsets[file] = offset
		offset += lineCounts[file]
	}
}

func (s *setImp) infoFor(p token.Pos) (int, string, int) {
	s.finish()
	if p <= token.NoPos {
		return 0, ``, 0
	}

	fsp := s.fs.Position(p)
	file, line := s.relPath(fsp.Filename), fsp.Line
	offset := s.offsets[file] + line - 1
	return offset, file, line
}

func (s *setImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	s.finish()
	m := jsonify.NewMap()
	files := utils.SortedKeys(s.offsets)
	for _, file := range files {
		offset := s.offsets[file]
		m.Add(ctx, strconv.Itoa(offset), file)
	}
	return m
}