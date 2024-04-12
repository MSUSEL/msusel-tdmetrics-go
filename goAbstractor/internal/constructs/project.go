package constructs

import "github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/jsonify"

type Project struct {
	Packages []*Package
}

func (p *Project) ToJson(ctx jsonify.Context) jsonify.Datum {
	return jsonify.NewMap().
		Add(ctx, `language`, `go`).
		AddNonZero(ctx, `packages`, p.Packages)
}
