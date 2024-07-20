package logger

import "fmt"

type outLog struct{}

func (outLog) Log(format string, args ...any) {
	fmt.Printf(format+"\n", args...)
}
