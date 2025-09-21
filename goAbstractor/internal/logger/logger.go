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

const indentText = `|  `
const showGroups = true

type LogEntry struct {
	Indent  int
	Group   string
	Message string
}

// New creates a new logger that writes to the standard out.
func New() *Logger {
	return &Logger{out: writerToFunc(os.Stdout)}
}

// NewWriter creates a logger that writes logs to the given writer.
//
// If the writer returns an error on any log, that error will be panicked
// from the call to create the log.
func NewWriter(out io.Writer) *Logger {
	if utils.IsNil(out) {
		return nil
	}
	return &Logger{out: writerToFunc(out)}
}

// NewFunc will create a logger that calls the given
// function to handle logging a message.
func NewFunc(handle func(entry LogEntry)) *Logger {
	if utils.IsNil(handle) {
		return nil
	}
	return &Logger{out: handle}
}

// Null is a nil logger that is always disabled.
// This can be used to prevent logging.
func Null() *Logger {
	return nil
}

// Logger is a tool for optionally logging messages.
type Logger struct {
	out      func(entry LogEntry)
	indent   int
	show     map[string]struct{}
	curGroup string
	disabled bool
}

func simpleToFunc(out func(string)) func(entry LogEntry) {
	return func(entry LogEntry) {
		prefix := "\n" + strings.Repeat(indentText, entry.Indent)
		if showGroups && len(entry.Group) > 0 {
			prefix += `[` + entry.Group + `] `
		}
		for line := range strings.SplitSeq(entry.Message, "\n") {
			out(prefix + line)
		}
	}
}

func writerToFunc(out io.Writer) func(entry LogEntry) {
	return simpleToFunc(func(text string) {
		if _, err := out.Write([]byte(text)); err != nil {
			panic(terror.New(`failed to write to a log`, err))
		}
	})
}

func (log *Logger) copy() *Logger {
	c := *log
	c.show = maps.Clone(log.show)
	return &c
}

func (log *Logger) write(msg string) *Logger {
	if log != nil && !log.disabled {
		log.out(LogEntry{
			Indent:  log.indent,
			Group:   log.curGroup,
			Message: msg,
		})
	}
	return log
}

// Log will write a log if this logger is enabled based on visible groups.
func (log *Logger) Log(args ...any) *Logger {
	return log.write(fmt.Sprint(args...))
}

// Logf will write a log if this logger is enabled based on visible groups.
func (log *Logger) Logf(format string, args ...any) *Logger {
	return log.write(fmt.Sprintf(format, args...))
}

// Indent will add to the current indent for the message.
// The returned logger will be further indented.
// If the logger is disabled, this will have no effect.
func (log *Logger) Indent() *Logger {
	if log == nil || log.disabled {
		return log
	}
	c := log.copy()
	c.indent++
	return c
}

// Disabled indicates that the current group is disabled.
func (log *Logger) Disabled() bool {
	return log == nil || log.disabled
}

func (log *Logger) updateEnable() {
	if len(log.curGroup) <= 0 {
		log.disabled = false
		return
	}
	_, has := log.show[log.curGroup]
	log.disabled = !has
}

// Group indicates that all the logs to the returned logger
// will be part of this group until another group is called.
// Only shown groups will be logged.
// If the group is an empty name, then it is shown.
func (log *Logger) Group(name string) *Logger {
	if log == nil {
		return nil
	}
	c := log.copy()
	c.curGroup = name
	c.updateEnable()
	return c
}

// Show indicates which groups will be logged. The receiver is not
// modified, the returned logger will have the shown groups in it.
// This will update if the Logger is disabled or not.
func (log *Logger) Show(groups ...string) *Logger {
	if log == nil {
		return nil
	}
	c := log.copy()
	if c.show == nil {
		c.show = make(map[string]struct{})
	}
	for _, group := range groups {
		if len(group) > 0 {
			c.show[group] = struct{}{}
		}
	}
	c.updateEnable()
	return c
}

// Hide indicates which groups will not be logged. The receiver is not
// modified, the returned logger will have the still shown groups in it.
// This will update if the Logger is disabled or not.
func (log *Logger) Hide(groups ...string) *Logger {
	if log == nil {
		return nil
	}
	c := log.copy()
	for _, group := range groups {
		delete(c.show, group)
	}
	c.updateEnable()
	return c
}
