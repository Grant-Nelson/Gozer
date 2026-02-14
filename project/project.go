package project

import (
	"go/token"
	"slices"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// Project represents all the information for a build of an application,
// including builds for running tests, benchmarks, or examples.
type Project struct {

	// FileSet is the file information for the package.
	FileSet *token.FileSet

	// Roots is the list of all root packages in the order given when
	// the project was created.
	Roots []*Package

	// AllPackages is the list of all packages in dependencies-first order,
	// meaning all dependencies for a package will be before the package
	// depending on those packages. All the root packages will be at or near
	// the end of the list.
	AllPackages []*Package

	// PackageMap is the set of all packages keyed by the package's import path.
	PackageMap map[string]*Package
}

func New(fSet *token.FileSet, roots []*packages.Package) *Project {
	proj := &Project{
		FileSet: fSet,
	}

	// Collect and prepare all the packages.
	pkgMap := map[string]*Package{}
	allBasePkgs := slices.Collect(packages.Postorder(roots))
	allPkgs := make([]*Package, len(allBasePkgs))
	for i, basePkg := range allBasePkgs {
		pkg := &Package{
			Project: proj,
			Ast:     basePkg,
		}
		allPkgs[i] = pkg
		pkgMap[basePkg.PkgPath] = pkg
	}
	proj.AllPackages = allPkgs
	proj.PackageMap = pkgMap

	// Get set of root packages.
	rootPkgs := make([]*Package, len(roots))
	for i, root := range roots {
		pkg, found := pkgMap[root.PkgPath]
		if !found {
			panic(faults.New(`failed to find root package in set of packages`).
				With(`path`, root.PkgPath))
		}
		rootPkgs[i] = pkg
	}
	proj.Roots = rootPkgs
	return proj
}
