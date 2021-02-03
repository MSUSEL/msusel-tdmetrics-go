package filter

import (
	"regexp"
	"strings"
)

// Matcher is a handler for comparing if a filename matches some criteria.
type Matcher func(filename string) bool

// TestFile returns true if the given filename is a test file.
func TestFile(filename string) bool {
	return strings.HasSuffix(filename, `_test.go`)
}

// GitFile returns true if the given filename is in a git path.
func GitFile(filename string) bool {
	return strings.Contains(filename, `/.git/`)
}

// VendorFile returns true if the given filename is in a vendor path.
func VendorFile(filename string) bool {
	return strings.Contains(filename, `/vendor/`)
}

// GoFile returns true if the given filename is a *.go file.
func GoFile(filename string) bool {
	return strings.HasSuffix(filename, `.go`)
}

// Not returns opposite result of the given matcher
func Not(matcher Matcher) Matcher {
	return func(filename string) bool {
		return !matcher(filename)
	}
}

// And returns true if all the matchers return true.
func And(matchers ...Matcher) Matcher {
	return func(filename string) bool {
		for _, matcher := range matchers {
			if !matcher(filename) {
				return false
			}
		}
		return len(matchers) > 0
	}
}

// Or returns true if any of the matchers return true.
func Or(matchers ...Matcher) Matcher {
	return func(filename string) bool {
		for _, matcher := range matchers {
			if matcher(filename) {
				return true
			}
		}
		return false
	}
}

// FileMatches returns true if the given filename is matched by the given regular expression.
func FileMatches(expression string) Matcher {
	re := regexp.MustCompile(expression)
	return func(filename string) bool {
		return re.Match([]byte(filename))
	}
}
