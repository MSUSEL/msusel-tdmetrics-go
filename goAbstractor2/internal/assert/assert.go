package assert

import (
	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func Unimplemented() {
	panic(terror.New(`Unimplemented`, nil))
}

func NotNil(v any, msg string) {
	if utils.IsNil(v) {
		panic(terror.New(`Value may not be nil`, nil).
			With(`message`, msg))
	}
}
