package main

import (
	"fmt"
	"os"

	"./internal/filter"
	"./internal/reader"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`Must provide only the path to the base folder of the go code.`)
		fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
		return
	}
	basePath := os.Args[1]

	reader.New().
		SetBasePath(basePath).
		AddDir(basePath).
		FilterFilenames(filter.Default()).
		Read()
}
