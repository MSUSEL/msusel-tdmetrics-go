package locs

import (
	"fmt"
	"go/token"

	"github.com/Snow-Gremlin/goToolbox/utils"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Loc interface {
	isLoc()

	Flag()
	Pos() token.Pos
	Info() (int, string, int)
	String() string
}

type locImp struct {
	s Set
	p token.Pos
}

func newLoc(s Set, p token.Pos) Loc {
	return &locImp{s: s, p: p}
}

func NoLoc() Loc {
	return newLoc(nil, token.NoPos)
}

func (c *locImp) isLoc() {}

func (c *locImp) Flag() {
	if !utils.IsNil(c.s) {
		c.s.flag(c.p)
	}
}

func (c *locImp) Pos() token.Pos {
	return c.p
}

func (c *locImp) Info() (int, string, int) {
	if utils.IsNil(c.s) {
		return 0, ``, 0
	}
	return c.s.infoFor(c.p)
}

func (c *locImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	offset, file, line := c.Info()

	if ctx.IsDebugFullLocIncluded() {
		return jsonify.NewMap().
			AddNonZero(ctx, `offset`, offset).
			AddNonZero(ctx, `file`, file).
			AddNonZero(ctx, `line`, line)
	}

	return jsonify.New(ctx, offset)
}

func (c *locImp) String() string {
	_, file, line := c.Info()
	return fmt.Sprintf(`%s:%d`, file, line)
}
