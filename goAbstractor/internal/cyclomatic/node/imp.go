package node

import (
	"fmt"
	"go/token"
	"strings"
)

type nodeImp struct {
	kind  string
	pos   token.Pos
	nexts []Node
}

func New(kind string, pos token.Pos) Node {
	return &nodeImp{
		kind: kind,
		pos:  pos,
	}
}

func (n *nodeImp) AddNext(next Node) {
	for _, prior := range n.nexts {
		if prior == next {
			return
		}
	}
	n.nexts = append(n.nexts, next)
}

func (n *nodeImp) Nexts() []Node {
	return n.nexts
}

func (n *nodeImp) Pos() token.Pos {
	return n.pos
}

func (n *nodeImp) String() string {
	return fmt.Sprintf(`%s_%d`, n.kind, n.pos)
}

func graphFormat(n Node, buf *strings.Builder, touched map[string]bool, first, rest string) {
	name := n.String()
	if touched[name] {
		buf.WriteString(fmt.Sprintf("%s──<%s>\n", first, name))
		return
	}
	touched[name] = true

	nexts := n.Nexts()
	max := len(nexts) - 1
	if max < 0 {
		buf.WriteString(fmt.Sprintf("%s──[%s]\n", first, name))
		return
	}

	buf.WriteString(fmt.Sprintf("%s┬─[%s]\n", first, name))
	for i, next := range nexts {
		if i == max {
			graphFormat(next, buf, touched, rest+`└─`, rest+`  `)
		} else {
			graphFormat(next, buf, touched, rest+`├─`, rest+`│ `)
		}
	}
}

func mermaidFormat(n Node, buf *strings.Builder, touched map[string]bool) {
	name := n.String()
	if touched[name] {
		return
	}
	touched[name] = true

	nexts := n.Nexts()
	buf.WriteString(fmt.Sprintf("   %s\n", name))
	for _, next := range nexts {
		buf.WriteString(fmt.Sprintf("   %s --> %s\n", name, next.String()))
	}
	for _, next := range nexts {
		mermaidFormat(next, buf, touched)
	}
}
