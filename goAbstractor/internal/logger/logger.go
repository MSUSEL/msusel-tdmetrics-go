package logger

import (
	"fmt"
	"io"
	"maps"
	"os"
	"strings"

	"github.com/Snow-Gremlin/goToolbox/terrors/terror"
	"github.com/Snow-Gremlin/goToolbox/utils"
)

func New() *Logger {
	return &Logger{out: writerToFunc(os.Stdout)}
}

func NewWriter(out io.Writer) *Logger {
	if utils.IsNil(out) {
		return nil
	}
	return &Logger{out: writerToFunc(out)}
}

func NewFunc(handle func(args ...any)) *Logger {
	if utils.IsNil(handle) {
		return nil
	}
	return &Logger{out: func(text string) { handle(text) }}
}

func Null() *Logger {
	return nil
}

type Logger struct {
	out    func(string)
	prefix string
	show   map[string]struct{}
}

func writerToFunc(out io.Writer) func(string) {
	return func(text string) {
		if _, err := out.Write([]byte(text)); err != nil {
			panic(terror.New(`failed to write to a log`, err))
		}
	}
}

func (log *Logger) copy() *Logger {
	return &Logger{
		out:    log.out,
		prefix: log.prefix,
		show:   maps.Clone(log.show),
	}
}

func (log *Logger) write(text string) *Logger {
	if log == nil {
		return nil
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		log.out("\n" + log.prefix + line)
	}
	return log
}

func (log *Logger) Log(args ...any) *Logger {
	return log.write(fmt.Sprint(args...))
}

func (log *Logger) Logf(format string, args ...any) *Logger {
	return log.write(fmt.Sprintf(format, args...))
}

// Prefix will add the given prefix to the end of the prior prefixes such that
// any logs to the returned logger will have the cumulated prior prefixes.
func (log *Logger) Prefix(prefix string) *Logger {
	if log == nil {
		return log
	}
	c := log.copy()
	c.prefix += prefix
	return c
}

// Group indicates that all the logs to the returned logger
// will be part of this group. Only shown groups will be logged.
func (log *Logger) Group(name string) *Logger {
	if log == nil {
		return nil
	}
	if _, has := log.show[name]; has {
		return log
	}
	return nil
}

// Show indicates which groups will be logged. The receiver is not
// modified, the returned logger will have the shown groups in it.
func (log *Logger) Show(groups ...string) *Logger {
	if log == nil {
		return nil
	}
	c := log.copy()
	if c.show == nil {
		c.show = make(map[string]struct{})
	}
	for _, group := range groups {
		c.show[group] = struct{}{}
	}
	return c
}
