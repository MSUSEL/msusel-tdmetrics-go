package utils

import (
	"sort"
	"strings"
)

// StringSet is a simple set of strings.
type StringSet struct {
	values map[string]bool
}

// NewStringSet creates a new set of strings.
func NewStringSet() *StringSet {
	return &StringSet{
		values: map[string]bool{},
	}
}

// Add will add all the given values into the set.
func (s *StringSet) Add(values ...string) {
	for _, value := range values {
		s.values[value] = true
	}
}

// Has will check that all the given values exist in the set.
// False will be returned if any one value is missing.
func (s *StringSet) Has(values ...string) bool {
	for _, value := range values {
		if !s.values[value] {
			return false
		}
	}
	return true
}

// Values will get the set of strings as a sorted slice.
func (s *StringSet) Values() []string {
	result := make([]string, len(s.values))
	i := 0
	for value := range s.values {
		result[i] = value
		i++
	}
	sort.Strings(result)
	return result
}

// String gets the debug string listing all the values in the set.
func (s *StringSet) String() string {
	if s == nil || len(s.values) <= 0 {
		return `<empty>`
	}
	return strings.Join(s.Values(), `, `)
}
