package project

import (
	"cmp"
	"go/token"
	"iter"
	"slices"

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
		FileSet:    fSet,
		PackageMap: map[string]*Package{},
	}
	for basePkg := range packages.Postorder(roots) {
		proj.insertPackage(basePkg)
	}
	for _, root := range roots {
		proj.assignRootPackage(root)
	}
	sortPackagesByDepth(proj.AllPackages)
	return proj
}

// insertPackage adds a new package to the project.
// This assumes all imports of this package have already been added
// and that the current package has not been added yet.
func (proj *Project) insertPackage(basePkg *packages.Package) {
	basePkg.Fset = proj.FileSet

	pkg := &Package{
		State: buildState.Listed,
		Ast:   basePkg,
	}
	proj.AllPackages = append(proj.AllPackages, pkg)
	proj.PackageMap[basePkg.PkgPath] = pkg

	depth := 0
	for pkgPath := range basePkg.Imports {
		depth = max(depth, proj.PackageMap[pkgPath].Depth+1)
	}
	pkg.Depth = depth
}

// assignRootPackage set an existing package as a root package.
// This assumes that the package has already been added.
func (proj *Project) assignRootPackage(root *packages.Package) {
	rootPkg := proj.PackageMap[root.PkgPath]
	proj.Roots = append(proj.Roots, rootPkg)
}

// sortPackagesByDepth performs stable sort of the given packages
// based on the depth of the package in the dependency tree.
func sortPackagesByDepth(pkgs []*Package) {
	slices.SortStableFunc(pkgs, func(a, b *Package) int {
		return cmp.Compare(a.Depth, b.Depth)
	})
}

func (proj *Project) UnfinishedPackages() iter.Seq[*Package] {
	return func(yield func(*Package) bool) {
		for _, pkg := range proj.AllPackages {
			if pkg.State != buildState.Finished && !yield(pkg) {
				return
			}
		}
	}
}
