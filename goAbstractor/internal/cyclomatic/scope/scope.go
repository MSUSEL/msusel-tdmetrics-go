package scope

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/cyclomatic/node"
)

const (
	// Enter is the start of a function.
	Enter = `$enter`

	// Exit is the stop of a function.
	// This is the node that `return` will connect to.
	// This will be changed by defers.
	Exit = `$exit`

	// Begin is the start of a block, like a loop.
	// This is the node that `continue` will connect to.
	Begin = `$begin`

	// End is the end of a block, like a loop.
	// This is the node that `break` will connect to.
	End = `$end`
)

type (
	Scope interface {
		Push() Scope
		Set(tag string, n node.Node)
		Get(tag string) node.Node
	}

	scopeImp struct {
		prior *scopeImp
		tags  map[string]node.Node
	}
)

func New() Scope {
	return &scopeImp{}
}

func (s *scopeImp) Push() Scope {
	return &scopeImp{prior: s}
}

func (s *scopeImp) Set(tag string, n node.Node) {
	if s.tags == nil {
		s.tags = map[string]node.Node{}
	}
	s.tags[tag] = n
}

func (s *scopeImp) Get(tag string) node.Node {
	for ; s != nil; s = s.prior {
		if n, has := s.tags[tag]; has {
			return n
		}
	}
	panic(fmt.Errorf(`error getting %s label in scope`, tag))
}

func (s *scopeImp) String() string {
	combo := dictionary.New[string, node.Node]()
	for ; s != nil; s = s.prior {
		combo.AddMapIfNotSet(s.tags)
	}
	return combo.String()
}
