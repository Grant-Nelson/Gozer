package astTools

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"sync"

	"golang.org/x/tools/go/packages"
)

var wd = sync.OnceValue(func() string {
	dir, err := filepath.Abs(`./`)
	if err != nil {
		panic(fmt.Errorf(`Error getting working directory: %w`, err))
	}
	return dir
})

func workingPath(path string) string {
	return filepath.Join(wd(), path)
}

func ParsePackage(inputFiles, extraFiles map[string]string) (*packages.Package, error) {
	patterns := make([]string, 0, len(inputFiles))
	overlay := make(map[string][]byte, len(inputFiles)+len(extraFiles))
	for name, src := range inputFiles {
		patterns = append(patterns, name)
		overlay[workingPath(name)] = []byte(src)
	}
	slices.Sort(patterns)
	for name, src := range extraFiles {
		overlay[workingPath(name)] = []byte(src)
	}

	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName |
			packages.NeedFiles |
			packages.NeedImports |
			packages.NeedDeps |
			packages.NeedTypes |
			packages.NeedTypesInfo |
			packages.NeedSyntax,
		Dir:     wd(),
		Overlay: overlay,
	}, patterns...)
	if err != nil {
		return nil, err
	}

	pkgErrors := []error{}
	packages.Visit(pkgs, nil, func(pkg *packages.Package) {
		for _, err := range pkg.Errors {
			pkgErrors = append(pkgErrors, err)
		}
	})
	if len(pkgErrors) > 0 {
		return nil, errors.Join(pkgErrors...)
	}
	if len(pkgs) != 1 {
		return nil, fmt.Errorf(`Expected exactly one root package but got %d`, len(pkgs))
	}
	return pkgs[0], nil
}
