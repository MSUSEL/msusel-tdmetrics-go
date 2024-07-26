package jsonify

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors"
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

var (
	seekKeyValuePattern = utils.LazyRegex(`^\s*(\w+)\s*=(.+)\s*$`)
	seekRangePattern    = utils.LazyRegex(`^\s*(\d+)\s*..\s*(\d+)\s*$`)
)

type seeker struct {
	path []any
	step int
}

func newSeeker(path []any) *seeker {
	return &seeker{path: path}
}

func (s *seeker) next() *seeker {
	return &seeker{path: s.path, step: s.step + 1}
}

func (s *seeker) done() bool {
	return s.step >= len(s.path)
}

func (s *seeker) fail(msg string, errs ...error) terrors.TError {
	return terror.New(msg, errs...).
		With(`step`, s.step).
		With(`path`, s.path)
}

func (s *seeker) raw() any {
	return s.path[s.step]
}

func (s *seeker) asInt() (int, bool) {
	switch cur := s.raw().(type) {
	case int:
		return cur, true
	case string:
		if v, err := strconv.ParseInt(cur, 0, 0); err == nil {
			return int(v), true
		}
	}
	return 0, false
}

func (s *seeker) asString() string {
	return fmt.Sprint(s.raw())
}

func (s *seeker) asIndex(length int) (int, bool) {
	index, ok := s.asInt()
	if !ok {
		return 0, false
	}
	if index < 0 || int(index) >= length {
		panic(s.fail(`index is out-of-bounds`).
			With(`index`, index).
			With(`length`, length))
	}
	return int(index), true
}

func (s *seeker) asRange(length int) (int, int, bool) {
	cur := s.asString()
	parts := seekRangePattern().FindAllStringSubmatch(cur, -1)
	if len(parts) != 1 || len(parts[0]) != 3 {
		return 0, 0, false
	}

	start, err := strconv.ParseInt(parts[0][1], 0, 0)
	if err != nil {
		return 0, 0, false
	}
	if start < 0 || int(start) >= length {
		panic(s.fail(`start index of a range is out-of-bounds`).
			With(`start`, start).
			With(`length`, length))
	}

	end, err := strconv.ParseInt(parts[0][2], 0, 0)
	if err != nil {
		return 0, 0, false
	}
	if end < start || int(end) >= length {
		panic(s.fail(`end index of a range is out-of-bounds`).
			With(`start`, start).
			With(`end`, end).
			With(`length`, length))
	}

	return int(start), int(end), true
}

func (s *seeker) getSelector(value string) (bool, func(string) bool) {
	value = strings.TrimSpace(value)

	var regexMatch bool
	value, regexMatch = strings.CutPrefix(value, `~`)

	if strings.HasPrefix(value, `"`) {
		var err error
		value, err = strconv.Unquote(value)
		if err != nil {
			panic(s.fail(`must have an unquotable string if quoted`, err).
				With(`value`, value))
		}
	}

	if !regexMatch {
		return true, func(v string) bool { return value == v }
	}

	r, err := regexp.Compile(value)
	if err != nil {
		panic(s.fail(`invalid regular expression following a '~'`, err).
			With(`value`, value))
	}
	return false, r.MatchString
}

func (s *seeker) asSelector() (bool, func(string) bool) {
	return s.getSelector(s.asString())
}

func (s *seeker) asKeyValue() (string, bool, func(string) bool, bool) {
	cur := s.asString()
	parts := seekKeyValuePattern().FindAllStringSubmatch(cur, -1)
	if len(parts) != 1 || len(parts[0]) != 3 {
		return ``, false, nil, false
	}

	single, selector := s.getSelector(parts[0][2])
	return parts[0][1], single, selector, true
}
