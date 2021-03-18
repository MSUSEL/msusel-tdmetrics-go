package main

import (
	"os"
	"path"

	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/filter"
	"github.com/MSUSEL/msusel-tdmetrics-go/exp001/internal/reader"
)

func main() {
	// if len(os.Args) != 2 {
	// 	fmt.Println(`Must provide only the path to the base folder of the go code.`)
	// 	fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
	// 	return
	// }
	// basePath := os.Args[1]

	basePath := path.Join(os.Getenv("GOPATH"), `src/github.com/MSUSEL/msusel-tdmetrics-go/testdata/dat002`)

	proj := reader.New().
		SetBasePath(basePath).
		AddDirRecursively(basePath).
		FilterFilenames(filter.Default()).
		Read()

	proj.PrintParticipation()
}
