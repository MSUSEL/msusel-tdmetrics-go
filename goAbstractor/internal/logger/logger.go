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
	return NewFor(os.Stdout)
}

func NewFor(out io.Writer) *Logger {
	return &Logger{out: out}
}

func Null() *Logger {
	return nil
}

type Logger struct {
	out    io.Writer
	prefix string
	show   map[string]struct{}
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
	if utils.IsNil(log.out) {
		log.out = os.Stdout
	}
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		if _, err := log.out.Write([]byte("\n" + log.prefix + line)); err != nil {
			panic(terror.New(`failed to write to a log`, err))
		}
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
