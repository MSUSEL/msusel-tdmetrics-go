package reader

import (
	"errors"
	"fmt"

	"golang.org/x/tools/go/packages"
)

// Read reads a project and all its packages and files.
func Read(config *Config) (ps []*packages.Package, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = recoverError(r)
		}
	}()

	cfg := config.toParseConfig()
	ps, err = packages.Load(cfg, config.Path)
	if err != nil {
		return nil, err
	}
	if err = allPackageErrors(ps); err != nil {
		return nil, err
	}
	return ps, nil
}

func recoverError(r any) error {
	switch r2 := r.(type) {
	case error:
		return r2
	case string:
		return errors.New(r2)
	case fmt.Stringer:
		return errors.New(r2.String())
	default:
		return fmt.Errorf(`error: %v`, r2)
	}
}

func allPackageErrors(ps []*packages.Package) error {
	errs := []error{}
	packages.Visit(ps, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			errs = append(errs, err)
		}
	})

	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errors.Join(errs...)
	}
}
