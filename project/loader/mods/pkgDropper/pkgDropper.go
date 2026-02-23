package pkgDropper

import (
	"maps"
	"path/filepath"
	"slices"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"golang.org/x/tools/go/packages"
)

type Config struct {
	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup

	// PkgPathPatterns is the list of package path patterns for packages to drop.
	// These patterns use `filepath.Match` to find packages to drop.
	PkgPathPatterns []string
}

type PkgDropper struct {
	errGroup        *faults.ErrGroup
	pkgPathPatterns []string
}

var _ mods.ModFactory = (*PkgDropper)(nil)

func New(cfg *Config) *PkgDropper {
	return &PkgDropper{
		errGroup:        cfg.ErrGroup,
		pkgPathPatterns: cfg.PkgPathPatterns,
	}
}

func (d *PkgDropper) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {
	return true, nil, nil
}

func (d *PkgDropper) pathMatcher(path string) bool {
	for _, pattern := range d.pkgPathPatterns {
		matched, err := filepath.Match(pattern, path)
		if err != nil {
			panic(err) // panic error so it is picked up by recover.
		}
		if matched {
			return true
		}
	}
	return false
}

func (d *PkgDropper) pkgMatcher(pkg *project.Package) bool {
	return d.pathMatcher(pkg.PkgPath())
}

func (d *PkgDropper) pkgMatcherWithKey(key string, pkg *project.Package) bool {
	return d.pathMatcher(pkg.PkgPath())
}

func (d *PkgDropper) basePkgMatcherWithKey(key string, pkg *packages.Package) bool {
	return d.pathMatcher(pkg.PkgPath)
}

func (d *PkgDropper) StartProject(proj *project.Project) (err error) {
	defer d.errGroup.Recover(&err)
	proj.AllPackages = slices.DeleteFunc(proj.AllPackages, d.pkgMatcher)
	proj.Roots = slices.DeleteFunc(proj.Roots, d.pkgMatcher)
	maps.DeleteFunc(proj.PackageMap, d.pkgMatcherWithKey)
	for _, p := range proj.AllPackages {
		maps.DeleteFunc(p.Ast.Imports, d.basePkgMatcherWithKey)
	}
	return d.errGroup.AnyOrNil()
}
