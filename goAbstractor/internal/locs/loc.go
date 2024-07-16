package locs

import (
	"go/token"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type (
	Loc interface {
		_loc()
		Flag()
		index() int
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

func (c *locImp) _loc() {}
func (c *locImp) Flag() { c.s.flag(c.p) }

func (c *locImp) index() int               { return c.s.indexFor(c.p) }
func (c *locImp) info() (int, string, int) { return c.s.infoFor(c.p) }

func (c *locImp) ToJson(ctx *jsonify.Context) jsonify.Datum {
	if ctx.IsFullLocationShown() {
		index, file, line := c.info()
		return jsonify.NewMap().
			AddNonZero(ctx, `index`, index).
			AddNonZero(ctx, `file`, file).
			AddNonZero(ctx, `line`, line)
	}

	return jsonify.New(ctx, c.index())
}
