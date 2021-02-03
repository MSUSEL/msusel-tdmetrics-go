package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`Must provide only the path to the base folder of the go code.`)
		fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
		return
	}
	basePath := os.Args[1]

	//startTime := time.Now()
	fmt.Println("Path:", basePath)
	//fmt.Printf("Done (%d files in %s)\n", len(p.Filenames()), time.Since(startTime).String())
}
