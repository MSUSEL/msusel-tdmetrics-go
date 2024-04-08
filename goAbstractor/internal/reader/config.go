package reader

import (
	"context"
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Config is the read and parse configuration.
type Config struct {

	// Verbose indicates that the parse logs
	// should be written to the standard out.
	Verbose bool

	// Path is the path to the main package or primary package.
	// The path should contain the mod file.
	// The path follows the standard patter for go tools.
	Path string

	// Context is the optional context to cancel a build with.
	Context context.Context

	// BuildFlags are the optional build flags to build with.
	// Example: // +build tag_name
	BuildFlags []string
}

func (c Config) toParseConfig() *packages.Config {
	const allNeeds = packages.NeedName |
		packages.NeedFiles |
		packages.NeedCompiledGoFiles |
		packages.NeedImports |
		packages.NeedDeps |
		packages.NeedExportFile |
		packages.NeedTypes |
		packages.NeedSyntax |
		packages.NeedTypesInfo |
		packages.NeedModule

	cfg := &packages.Config{
		Dir:        c.Path,
		BuildFlags: c.BuildFlags,
		Context:    c.Context,
		Mode:       allNeeds,
	}

	if c.Verbose {
		cfg.Logf = func(format string, args ...any) {
			_, err := fmt.Printf(format, args...)
			panic(err)
		}
	}

	return cfg
}
