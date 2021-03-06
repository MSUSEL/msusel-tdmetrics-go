package main

import (
	"fmt"
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

	basePath := path.Join(os.Getenv("GOPATH"), `src/github.com/MSUSEL/msusel-tdmetrics-go/testdata/dat001`)

	project := reader.New().
		SetBasePath(basePath).
		AddDirRecursively(basePath).
		FilterFilenames(filter.Default()).
		Read()

	for _, pkg := range project.Packages {
		fmt.Printf("Package: %q\n", pkg.Package.Path())
		fmt.Printf("Name:    %s\n", pkg.Package.Name())
		fmt.Println("===========================")
		for id, f := range pkg.Participation() {
			fmt.Println(id, "=>", f)
		}
		fmt.Println()
	}
}
