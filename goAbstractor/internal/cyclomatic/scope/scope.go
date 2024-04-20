package scope

import (
	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/cyclomatic/node"
)

const (
	Begin = `$begin`
	End   = `$end`

	EndIf = `$endIf`
)

type Scope interface {
	Push() Scope
	Set(tag string, n node.Node) Scope
	SetRange(begin, end node.Node) Scope
	Get(tag string) node.Node
}
