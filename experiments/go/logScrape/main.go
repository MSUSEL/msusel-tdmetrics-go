package main

import (
	"bytes"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"unicode/utf16"
	"unsafe"
)

const (
	toRepo  = `../../..`
	logFile = `javaAbstractor\tdd\commons-bcel.log`
	target  = `unimplemented addUsage:`
)

func main() {
	parts := slices.Collect(
		dedup(
			trimmed(
				contains(target,
					strings.Lines(
						readText(filepath.Join(toRepo, logFile)),
					),
				),
			),
		),
	)

	slices.Sort(parts)
	for _, part := range parts {
		fmt.Println(part)
	}
}

func cleanup(line string) string {
	line = cutFront(line, target)
	line = cutEnd(line, `@`)
	line = cutFront(line, `(`)
	line = cutEnd(line, `)`)
	line = cutFront(line, `spoon.support.reflect.code.`)
	line = cutEnd(line, `Impl`)
	return strings.TrimSpace(line)
}

func cutFront(s, sub string) string {
	if index := strings.Index(s, sub); index >= 0 {
		s = s[index+len(sub):]
	}
	return s
}

func cutEnd(s, sub string) string {
	if index := strings.Index(s, sub); index >= 0 {
		s = s[:index]
	}
	return s
}

func contains(sub string, lines iter.Seq[string]) iter.Seq[string] {
	return where(lines, func(s string) bool { return strings.Contains(s, sub) })
}

func trimmed(s iter.Seq[string]) iter.Seq[string] {
	return convert(s, cleanup)
}

func dedup[T comparable](s iter.Seq[T]) iter.Seq[T] {
	seen := map[T]bool{}
	return where(s, func(v T) bool {
		if seen[v] {
			return false
		}
		seen[v] = true
		return true
	})
}

func convert[T, U any](s iter.Seq[T], fn func(T) U) iter.Seq[U] {
	return func(yield func(U) bool) {
		for v := range s {
			if !yield(fn(v)) {
				return
			}
		}
	}
}

func where[T any](s iter.Seq[T], predicate func(T) bool) iter.Seq[T] {
	return func(yield func(T) bool) {
		for v := range s {
			if predicate(v) && !yield(v) {
				return
			}
		}
	}
}

func readText(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	if rest, has := bytes.CutPrefix(data, []byte{0xff, 0xfe}); has {
		dat16 := unsafe.Slice((*uint16)(unsafe.Pointer(unsafe.SliceData(rest))), len(rest)/2)
		return string(utf16.Decode(dat16))
	}
	return string(data)
}
