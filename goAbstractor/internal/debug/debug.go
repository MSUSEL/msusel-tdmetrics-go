package debug

import (
	"runtime/debug"
	"strings"
)

func StackCalls(offset, count int) string {
	lines := strings.Split(string(debug.Stack()), "\n")
	start := 5 + offset*2
	end := min(start+count*2, len(lines)-1)
	if end <= start {
		return `no-stack`
	}
	return strings.Join(lines[start:end], "\n")
}
