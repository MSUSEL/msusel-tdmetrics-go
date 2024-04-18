package node

import (
	"fmt"
	"go/token"
	"strings"
)

type (
	Node interface {
		AddNext(next Node)
		Nexts() []Node
		Pos() token.Pos
		String() string
		FullString() string
	}

	nodeImp struct {
		kind  string
		pos   token.Pos
		nexts []Node
	}
)

func New(kind string, pos token.Pos) Node {
	return &nodeImp{
		kind: kind,
		pos:  pos,
	}
}

func (n *nodeImp) AddNext(next Node) {
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

func (n *nodeImp) FullString() string {
	buf := &strings.Builder{}
	touched := map[string]bool{}
	format(n, buf, touched, `─`, ` `)
	return buf.String()
}

func format(n Node, buf *strings.Builder, touched map[string]bool, first, rest string) {
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
			format(next, buf, touched, rest+`└─`, rest+`  `)
		} else {
			format(next, buf, touched, rest+`├─`, rest+`│ `)
		}
	}
}
