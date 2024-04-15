package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

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

func (p *Package) String() string {
	return jsonify.ToString(p)
}

func (p *Package) FindTypeForReceiver(receiver string) *TypeDef {
	for _, t := range p.Types {
		if receiver == t.Name {
			return t
		}
	}
	return nil
}
