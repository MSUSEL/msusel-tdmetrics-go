package stringer

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/collections"
	"github.com/Snow-Gremlin/goToolbox/collections/set"
)

type Stringer interface {
	WriteString(text string) Stringer
	Write(args ...any) Stringer
	WriteList(open, sep, close string, list any) Stringer
	Reset() Stringer
	String() string
}

type Stringerable interface {
	ToStringer(s Stringer)
}

type stringerImp struct {
	buf    *strings.Builder
	marked collections.Set[Stringerable]
}

func New() Stringer {
	return &stringerImp{
		buf:    &strings.Builder{},
		marked: set.New[Stringerable](),
	}
}

func String(args ...any) string {
	return New().Write(args...).String()
}

func (s *stringerImp) mark(s2 Stringerable) bool {
	return s2 != nil && s.marked.Add(s2)
}

func (s *stringerImp) WriteString(text string) Stringer {
	if _, err := s.buf.WriteString(text); err != nil {
		panic(err)
	}
	return s
}

func (s *stringerImp) writeOne(arg any) {
	switch t := arg.(type) {
	case Stringerable:
		if s.mark(t) {
			t.ToStringer(s)
		} else {
			s.WriteString(`λ`)
		}
	case fmt.Stringer:
		s.WriteString(t.String())
	case string:
		s.WriteString(t)
	default:
		s.WriteString(fmt.Sprintf(`%v`, arg))
	}
}

func (s *stringerImp) Write(args ...any) Stringer {
	for _, arg := range args {
		s.writeOne(arg)
	}
	return s
}

func (s *stringerImp) WriteList(open, sep, close string, list any) Stringer {
	r := reflect.ValueOf(list)
	if r.Kind() != reflect.Slice && r.Kind() != reflect.Array {
		panic(fmt.Errorf("WriteList: expected slice or array, got %T", list))
	}

	count := r.Len()
	if count <= 0 {
		return s
	}

	s.WriteString(open)
	for i := range count {
		if i > 0 {
			s.WriteString(sep)
		}
		s.Write(r.Index(i).Interface())
	}
	s.WriteString(close)
	return s
}

func (s *stringerImp) Reset() Stringer {
	s.buf.Reset()
	return s
}

func (s *stringerImp) String() string {
	return s.buf.String()
}
