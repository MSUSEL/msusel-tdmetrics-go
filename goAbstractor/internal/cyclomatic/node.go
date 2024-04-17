package cyclomatic

import (
	"fmt"
	"go/token"
	"strings"
)

type node struct {
	kind  string
	pos   token.Pos
	nexts []*node
}

func newNode(kind string, pos token.Pos) *node {
	return &node{
		kind: kind,
		pos:  pos,
	}
}

func (n *node) addNext(next *node) {
	n.nexts = append(n.nexts, next)
}

func (n *node) format(buf *strings.Builder, touched map[string]bool, first, rest string) {
	name := n.String()
	if touched[name] {
		buf.WriteString(fmt.Sprintf("%s──<%s>\n", first, name))
		return
	}

	touched[name] = true
	if len(n.nexts) <= 0 {
		buf.WriteString(fmt.Sprintf("%s──[%s]\n", first, name))
		return
	}

	buf.WriteString(fmt.Sprintf("%s┬─[%s]\n", first, name))
	max := len(n.nexts) - 1
	for i, next := range n.nexts {
		if i == max {
			next.format(buf, touched, rest+`└─`, rest+`  `)
		} else {
			next.format(buf, touched, rest+`├─`, rest+`│ `)
		}
	}
}

func (n *node) String() string {
	return fmt.Sprintf(`%s_%d`, n.kind, n.pos)
}
