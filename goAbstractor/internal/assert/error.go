package assert

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

var idMatch = utils.LazyMatcher(`^\$?[a-zA-Z][_a-zA-Z0-9]*$`)

func ArgValidId(name, value string) {
	if !idMatch(value) {
		panic(terror.New(`argument must be a valid identifier`).
			With(`name`, name).
			With(`value`, value))
	}
}

func ArgNotNil(name string, value any) {
	if utils.IsNil(value) {
		panic(terror.New(`argument must not be nil`).
			With(`name`, name))
	}
}

func ArgNotEmpty(name, value string) {
	if len(value) <= 0 {
		panic(terror.New(`argument must not be an empty string`).
			With(`name`, name))
	}
}
