package assert

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

var idMatch = utils.LazyMatcher(`^\$?[_a-zA-Z][_a-zA-Z0-9]*(#[0-9]+)?$`)

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

func ArgNotEmpty(name string, value any) {
	if len, ok := utils.Length(value); !ok || len <= 0 {
		panic(terror.New(`argument must not be an empty string`).
			With(`name`, name))
	}
}

func ArgNoNils[T any, S ~[]T](name string, values S) {
	for _, v := range values {
		if utils.IsNil(v) {
			panic(terror.New(`slice may not contain a nil`).
				With(`name`, name))
		}
	}
}
