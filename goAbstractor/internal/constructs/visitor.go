package constructs

import "github.com/Snow-Gremlin/goToolbox/utils"

type Visitor func(value Visitable) bool

type Visitable interface {
	Visit(v Visitor)
}

func visitTest[T Visitable](v Visitor, value T) {
	if !utils.IsNil(value) && v(value) {
		value.Visit(v)
	}
}

func visitList[T Visitable, S ~[]T](v Visitor, values S) {
	for _, val := range values {
		visitTest(v, val)
	}
}

func visitMap[T Visitable, M ~map[string]T](v Visitor, values M) {
	for _, val := range values {
		visitTest(v, val)
	}
}
