package jsonify

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

var (
	seekWordMatcher     = utils.LazyMatcher(`^\w+$`)
	seekKeyValuePattern = utils.LazyRegex(`^\s*(\w+)\s*=(.+)\s*$`)
	seekRangePattern    = utils.LazyRegex(`^\s*(\d+)\s*..\s*(\d+)\s*$`)
)

type seeker struct {
	path  []any
	index int
}

func newSeeker(path []any) *seeker {
	return &seeker{path: path}
}

func (s *seeker) copy() *seeker {
	return &seeker{
		path:  s.path,
		index: s.index,
	}
}

func (s *seeker) done() bool {
	return s.index >= len(s.path)-1
}

func (s *seeker) current() any {
	return s.path[s.index]
}

func (s *seeker) unquoteValue(value string) string {
	if strings.HasPrefix(value, `"`) {
		var err error
		value, err = strconv.Unquote(value)
		if err != nil {
			panic(terror.New(`must have a value that is unquotable`, err).
				With(`value`, value).
				With(`path`, s.path).
				With(`step#`, s.index).
				With(`step`, s.current()))
		}
	}
	return value
}

func (s *seeker) nextStep(d Datum) Datum {
	s.index++
	switch dt := d.(type) {
	case *List:
		return s.StepList(dt)
	case *Map:
		return s.StepMap(dt)
	case *null:
		return s.StepNull(dt)
	default:
		return s.StepValue(dt)
	}
}

func (s *seeker) StepList(d *List) Datum {
	if s.done() {
		return d
	}

	if index, ok := s.current().(int); ok {
		return s.seekIndex(d, index)
	}

	if str, ok := s.current().(string); ok {
		if str == `#` {
			return newValue(len(d.data))
		}
		if seekKeyValuePattern().MatchString(str) {
			return s.seekKeyValue(d, str)
		}
		if seekRangePattern().MatchString(str) {
			return s.seekRange(d, str)
		}
		if index, err := strconv.ParseInt(str, 0, 0); err != nil {
			return s.seekIndex(d, int(index))
		}
	}

	panic(terror.New(`must have an index, range, or key/value`).
		With(`path`, s.path).
		With(`step#`, s.index).
		With(`step`, s.current()))
}

func (s *seeker) checkRange(d *List, index int) {
	count := len(d.data)
	if index < 0 || index >= count {
		panic(terror.New(`index is out of bounds`).
			With(`index`, index).
			With(`count`, count).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}
}

func (s *seeker) parseIndex(d *List, str string) int {
	index, err := strconv.ParseInt(str, 0, 0)
	if err != nil {
		panic(terror.New(`must have an int parsable index`, err).
			With(`index`, str).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}
	s.checkRange(d, int(index))
	return int(index)
}

func (s *seeker) seekIndex(d *List, index int) Datum {
	s.checkRange(d, index)
	return s.nextStep(d.data[index])
}

func (s *seeker) seekRange(d *List, rng string) Datum {
	parts := seekRangePattern().FindAllStringSubmatch(rng, -1)
	if len(parts) != 1 || len(parts[0]) != 3 {
		panic(terror.New(`invalid range in path`).
			With(`pattern`, `^(\d+)..(\d+)$`).
			With(`gotten`, parts).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}

	start := s.parseIndex(d, parts[0][1])
	end := s.parseIndex(d, parts[0][2])
	sub := NewList()
	for i := start; i <= end; i++ {
		d := s.copy().nextStep(d.data[i])
		sub.data = append(sub.data, d)
	}
	return sub
}

func (s *seeker) seekKeyValue(d *List, kv string) Datum {
	parts := seekKeyValuePattern().FindAllStringSubmatch(kv, -1)
	if len(parts) != 1 || len(parts[0]) != 3 {
		panic(terror.New(`invalid key/value in path`).
			With(`pattern`, `^(\w+)=(.+)$`).
			With(`gotten`, parts).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}

	key, value := parts[0][1], s.unquoteValue(parts[0][2])

	foundKey := false
	if seekWordMatcher(value) {
		for _, item := range d.data {
			if m, ok := item.(*Map); ok {
				if v := m.Get(key); !utils.IsNil(v) {
					foundKey = true
					if fmt.Sprint(v.RawValue()) == value {
						return s.nextStep(item)
					}
				}
			}
		}
	} else {
		re, err := regexp.Compile(value)
		if err != nil {
			panic(terror.New(`regex value for list of objects not compilable`, err).
				With(`key`, key).
				With(`path`, s.path).
				With(`step#`, s.index).
				With(`step`, s.current()))
		}

		l := NewList()
		for _, item := range d.data {
			if m, ok := item.(*Map); ok {
				if v := m.Get(key); !utils.IsNil(v) {
					foundKey = true
					if re.MatchString(fmt.Sprint(v.RawValue())) {
						e := s.nextStep(item)
						l.data = append(l.data, e)
					}
				}
			}
		}
		if len(l.data) > 0 {
			return l
		}
	}

	if foundKey {
		panic(terror.New(`no value found with the given key`).
			With(`key`, key).
			With(`value`, value).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}

	panic(terror.New(`no key found`).
		With(`key`, key).
		With(`path`, s.path).
		With(`step#`, s.index).
		With(`step`, s.current()))
}

func (s *seeker) StepMap(d *Map) Datum {
	if s.done() {
		return d
	}

	key := s.unquoteValue(fmt.Sprint(s.current()))
	if seekWordMatcher(key) {
		value, ok := d.data[key]
		if !ok {
			panic(terror.New(`key not found in map`).
				With(`key`, key).
				With(`path`, s.path).
				With(`step#`, s.index).
				With(`step`, s.current()))
		}
		return s.nextStep(value)
	}

	re, err := regexp.Compile(key)
	if err != nil {
		panic(terror.New(`regex key for map not compilable`, err).
			With(`key`, key).
			With(`path`, s.path).
			With(`step#`, s.index).
			With(`step`, s.current()))
	}

	l := NewList()
	keys := utils.SortedKeys(d.data)
	for _, key := range keys {
		if re.MatchString(key) {
			e := s.copy().nextStep(d.data[key])
			l.data = append(l.data, e)
		}
	}
	return l
}

func (s *seeker) StepNull(d *null) Datum {
	if s.done() {
		return d
	}

	panic(terror.New(`path continues from null`).
		With(`path`, s.path).
		With(`step#`, s.index).
		With(`step`, s.current()))
}

func (s *seeker) StepValue(d Datum) Datum {
	if s.done() {
		return d
	}

	panic(terror.New(`path continues from the value`).
		With(`data`, d.RawValue()).
		With(`path`, s.path).
		With(`step#`, s.index).
		With(`step`, s.current()))
}
