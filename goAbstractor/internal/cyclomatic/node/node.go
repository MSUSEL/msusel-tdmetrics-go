package node

import (
	"go/token"
	"strings"
)

type Node interface {
	AddNext(next Node)
	Nexts() []Node
	Pos() token.Pos
	String() string
}

func Graph(n Node) string {
	buf := &strings.Builder{}
	touched := map[string]bool{}
	graphFormat(n, buf, touched, `─`, ` `)
	return buf.String()
}

func Mermaid(n Node) string {
	buf := &strings.Builder{}
	touched := map[string]bool{}
	buf.WriteString("stateDiagram-v2\n")
	mermaidFormat(n, buf, touched)
	return buf.String()
}
