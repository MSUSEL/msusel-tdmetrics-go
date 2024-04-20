package scope

import (
	"fmt"

	"github.com/Snow-Gremlin/goToolbox/collections/dictionary"

	"github.com/MSUSEL/msusel-tdmetrics-go/goAbstractor/internal/cyclomatic/node"
)

type scopeImp struct {
	prior *scopeImp
	tags  map[string]node.Node
}

func New() Scope {
	return &scopeImp{}
}

func (s *scopeImp) Push() Scope {
	return &scopeImp{prior: s}
}

func (s *scopeImp) Set(tag string, n node.Node) Scope {
	if s.tags == nil {
		s.tags = map[string]node.Node{}
	}
	s.tags[tag] = n
	return s
}

func (s *scopeImp) SetRange(begin, end node.Node) Scope {
	return s.Set(Begin, begin).Set(End, end)
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
