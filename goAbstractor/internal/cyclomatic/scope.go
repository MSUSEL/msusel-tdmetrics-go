package cyclomatic

import "fmt"

const (
	enterTag = `$enter`
	exitTag  = `$exit`
	beginTag = `$begin`
	endTag   = `$end`
)

type scope struct {
	prior *scope
	tags  map[string]*node
}

func newScope() *scope {
	return &scope{}
}

func (s *scope) push() *scope {
	return &scope{prior: s}
}

func (s *scope) setTag(name string, n *node) {
	if s.tags == nil {
		s.tags = map[string]*node{}
	}
	s.tags[name] = n
}

func (s *scope) getTag(name string) *node {
	if s == nil {
		panic(fmt.Errorf(`error getting %s label in scope`, name))
	}
	if n, has := s.tags[name]; has {
		return n
	}
	return s.prior.getTag(name)
}
