package locs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

type (
	Loc interface {
		_loc()
		Flag()
		info() (int, string, int)
	}

	locImp struct {
		s Set
		p token.Pos
	}
)

func newLoc(s Set, p token.Pos) Loc {
	return &locImp{s: s, p: p}
}

func NoLoc() Loc {
	return newLoc(nil, token.NoPos)
}

func (c *locImp) _loc() {}

func (c *locImp) Flag() {
	if !utils.IsNil(c.s) {
		c.s.flag(c.p)
	}
}

func (c *locImp) info() (int, string, int) {
	if utils.IsNil(c.s) {
		return 0, ``, 0
	}
	return c.s.infoFor(c.p)
}

func (c *locImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	offset, file, line := c.info()

	if ctx.IsFullLocationShown() {
		return jsonify.NewMap().
			AddNonZero(ctx, `offset`, offset).
			AddNonZero(ctx, `file`, file).
			AddNonZero(ctx, `line`, line)
	}

	return jsonify.New(ctx, offset)
}
