package main

import (
	"fmt"
	"os"
	"time"

	"./internal/handlers"
	"./internal/parser"
	"./internal/processors"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println(`Must provide only the path to the base folder of the go code.`)
		fmt.Println(`Example: go run main.go /Users/<user name>/go/src/github.com/MSUSEL/msusel-tdmetrics-go`)
		return
	}
	basePath := os.Args[1]

	proc := processors.NewFuncUsage()
	p := parser.NewParser().
		//UpdateProgress(parser.MuteProgress).
		AddProcessor(proc.ProcessFunction()).
		AddHandler(handlers.JustFunctionParameterTypes)

	p.AddDirRecursively(basePath).
		FilterFilenames(parser.IsTestFile).
		FilterFilenames(parser.IsGitFile).
		FilterFilenames(parser.IsNotGoFile).
		FilterFilenames(parser.IsVendorFile)
	//fmt.Println(strings.Join(p.Filenames(), "\n"))

	startTime := time.Now()
	p.Start().Await()
	fmt.Printf("Done (%d files in %s)\n", len(p.Filenames()), time.Since(startTime).String())

	fmt.Println()
	fmt.Println(proc.String())
}
