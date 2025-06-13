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

func ArgIsNil(name string, value any) {
	if !utils.IsNil(value) {
		panic(terror.New(`argument must be nil`).
			With(`name`, name))
	}
}

func ArgNotNil(name string, value any) {
	if utils.IsNil(value) {
		panic(terror.New(`argument must not be nil`).
			With(`name`, name))
	}
}

func ArgIsEmpty(name string, value any) {
	if len, ok := utils.Length(value); !ok || len >= 0 {
		panic(terror.New(`argument must be empty`).
			With(`name`, name))
	}
}

func ArgsHaveSameLength(names string, value1, value2 any) {
	len1, ok := utils.Length(value1)
	if !ok {
		panic(terror.New(`first argument unable to get length`).
			With(`names`, names))
	}
	len2, ok := utils.Length(value2)
	if !ok {
		panic(terror.New(`second argument unable to get length`).
			With(`names`, names))
	}
	if len1 != len2 {
		panic(terror.New(`argument must be expected length`).
			With(`first count`, len1).
			With(`second count`, len2).
			With(`first`, value1).
			With(`second`, value2).
			With(`names`, names))
	}
}

func ArgNotEmpty(name string, value any) {
	if len, ok := utils.Length(value); !ok || len <= 0 {
		panic(terror.New(`argument must not be empty`).
			With(`name`, name))
	}
}

func AnyArgNotEmpty(names string, values ...any) {
	for _, value := range values {
		if len, ok := utils.Length(value); ok && len > 0 {
			return
		}
	}
	panic(terror.New(`at least one argument must not be empty`).
		With(`names`, names))
}

func ArgHasNoNils[T any, S ~[]T](name string, values S) {
	for _, v := range values {
		if utils.IsNil(v) {
			panic(terror.New(`slice may not contain a nil`).
				With(`name`, name))
		}
	}
}
