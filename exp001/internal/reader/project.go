package reader

import (
	"go/token"
)

// Project is the collection of compiled data for the project.
type Project struct {
	BasePath string
	FileSet  *token.FileSet
	Packages []*Package
}
