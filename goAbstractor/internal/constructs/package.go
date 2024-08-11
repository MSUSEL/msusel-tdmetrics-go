package constructs

import "golang.org/x/tools/go/packages"

type Package interface {
	Construct
	Source() *packages.Package
}
