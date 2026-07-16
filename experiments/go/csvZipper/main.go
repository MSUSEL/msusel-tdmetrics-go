package main

import (
	"errors"
	"fmt"
	"iter"
	"os"
	"path/filepath"
	"strings"
)

const (
	toRepo  = `../../..`
	logFile = `javaAbstractor\tdd\metrics_output\`
)

func main() {
	files := filter(walkDirFiles(filepath.Join(toRepo, logFile)))
	for file := range files {
		zipFile(file)
	}
}

func zipFile(path string) {
	// TODO: If this is needed, finish it.
	fmt.Println(path)
}

func filter(files iter.Seq[string]) iter.Seq[string] {
	return where(files, func(s string) bool {
		return strings.HasSuffix(s, `.csv`)
	})
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

func walkDirFiles(root string) iter.Seq[string] {
	var errYieldBreak = errors.New(`errYieldBreak`)
	return func(yield func(string) bool) {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() {
				if !yield(path) {
					return errYieldBreak
				}
			}
			return nil
		})
		if err != nil {
			if err == errYieldBreak {
				return
			}
			panic(err)
		}
	}
}
