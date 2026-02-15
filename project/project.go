package project

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
	"golang.org/x/tools/go/packages"
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

// New constructs a new project and initializes that project
// with the given package information.
// The current packages are expected to only have file names listed
// and no syntax nor types determined yet.
func New(fSet *token.FileSet, roots []*packages.Package) *Project {
	proj := &Project{
		FileSet: fSet,
	}

	// Collect and prepare all the packages.
	pkgMap := map[string]*Package{}
	allPkgs := []*Package{}
	for basePkg := range packages.Postorder(roots) {
		pkg := &Package{
			State: buildState.Listed,
			Ast:   basePkg,
		}
		allPkgs = append(allPkgs, pkg)
		pkgMap[basePkg.PkgPath] = pkg
	}
	proj.AllPackages = allPkgs
	proj.PackageMap = pkgMap

	// Get set of root packages from base package root.
	rootPkgs := make([]*Package, len(roots))
	for i, root := range roots {
		rootPkgs[i] = pkgMap[root.PkgPath]
	}
	proj.Roots = rootPkgs
	return proj
}
