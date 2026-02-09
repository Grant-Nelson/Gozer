package project

import (
	"fmt"
	"go/token"
	"slices"

	"github.com/Grant-Nelson/Gozer/avail/iterator"
	"golang.org/x/tools/go/packages"
)

type Project struct {
	FileSet     *token.FileSet
	Roots       []*Package
	AllPackages []*Package
}

func New(fSet *token.FileSet, roots []*packages.Package) *Project {
	allBasePkgs := slices.Collect(packages.Postorder(roots))
	allPkgs := make([]*Package, len(allBasePkgs))
	for i, basePkg := range allBasePkgs {
		allPkgs[i] = &Package{Package: basePkg}
	}

	rootPkgs := make([]*Package, len(roots))
	for i, root := range roots {
		pkg, found := iterator.Iterate(allPkgs...).
			Where(func(p *Package) bool { return p.Package == root }).
			First()
		if !found {
			panic(fmt.Errorf(`failed to find root package in set of packages`))
		}
		rootPkgs[i] = pkg
	}

	return &Project{
		FileSet:     fSet,
		Roots:       rootPkgs,
		AllPackages: allPkgs,
	}
}
