package constructs

import (
	"encoding/json"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"
)

type Package struct {
	Path    string
	Imports []string
	Types   []*TypeDef
	Values  []*ValueDef
	Methods []*Method
}

func (p *Package) ToJson(ctx *jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		AddNonZero(ctx, `path`, p.Path).
		AddNonZero(ctx, `imports`, p.Imports).
		AddNonZero(ctx, `types`, p.Types).
		AddNonZero(ctx, `values`, p.Values).
		AddNonZero(ctx, `methods`, p.Methods)
}

func (p *Package) MarshalJSON() ([]byte, error) {
	ctx := jsonify.NewContext()
	return json.Marshal(p.ToJson(ctx))
}
