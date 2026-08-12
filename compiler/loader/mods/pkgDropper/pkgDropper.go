package pkgDropper

import (
	"maps"
	"path/filepath"
	"slices"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/logger"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type Config struct {

	// Logger to log verbose messages with. Has no affect if verbose was false.
	Logger *logger.Logger

	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup

	// PkgPathPatterns is the list of package path patterns for packages to drop.
	// These patterns use `filepath.Match` to find packages to drop.
	PkgPathPatterns []string
}

type PkgDropper struct {
	logger          *logger.Logger
	errGroup        *faults.ErrGroup
	pkgPathPatterns []string
}

var _ mods.ModFactory = (*PkgDropper)(nil)

func New(cfg *Config) *PkgDropper {
	return &PkgDropper{
		logger:          cfg.Logger,
		errGroup:        cfg.ErrGroup,
		pkgPathPatterns: cfg.PkgPathPatterns,
	}
}

func (d *PkgDropper) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {
	return true, nil, nil
}

func (d *PkgDropper) StartProject(proj *project.Project) (err error) {
	defer d.errGroup.Recover(&err)

	var removed map[string]struct{}
	if d.logger != nil {
		removed = make(map[string]struct{}, len(proj.AllPackages))
		for _, p := range proj.AllPackages {
			removed[p.Ast.PkgPath] = struct{}{}
		}
	}

	if !removePackages(d.pkgPathPatterns, proj) {
		return d.errGroup.AnyOrNil()
	}
	removeUnreachable(proj)

	if d.logger != nil {
		for _, p := range proj.AllPackages {
			delete(removed, p.Ast.PkgPath)
		}
		removedPaths := slices.Collect(maps.Keys(removed))
		slices.Sort(removedPaths)
		d.logger.Printf("Dropped Packages:\n\t%s\n", strings.Join(removedPaths, "\n\t"))
	}

	return d.errGroup.AnyOrNil()
}

func pkgPathMatcher(pkgPathPatterns []string) func(path string) bool {
	return func(path string) bool {
		for _, pattern := range pkgPathPatterns {
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
}

func pkgMatcher(pathMatcher func(path string) bool) func(pkg *project.Package) bool {
	return func(pkg *project.Package) bool { return pathMatcher(pkg.PkgPath()) }
}

func pkgMatcherWithKey(pathMatcher func(path string) bool) func(key string, pkg *project.Package) bool {
	return func(key string, pkg *project.Package) bool { return pathMatcher(pkg.PkgPath()) }
}

func basePkgMatcherWithKey(pathMatcher func(path string) bool) func(key string, pkg *packages.Package) bool {
	return func(key string, pkg *packages.Package) bool { return pathMatcher(pkg.PkgPath) }
}

func removePackages(pkgPathPatterns []string, proj *project.Project) bool {
	if len(pkgPathPatterns) <= 0 {
		return false
	}
	pathMatcher := pkgPathMatcher(pkgPathPatterns)

	count := len(proj.AllPackages)
	proj.AllPackages = slices.DeleteFunc(proj.AllPackages, pkgMatcher(pathMatcher))
	if count == len(proj.AllPackages) {
		return false
	}

	proj.Roots = slices.DeleteFunc(proj.Roots, pkgMatcher(pathMatcher))
	maps.DeleteFunc(proj.PackageMap, pkgMatcherWithKey(pathMatcher))
	for _, p := range proj.AllPackages {
		maps.DeleteFunc(p.Ast.Imports, basePkgMatcherWithKey(pathMatcher))
	}
	return true
}

func removeUnreachable(proj *project.Project) bool {
	baseRoots := make([]*packages.Package, len(proj.Roots))
	for i, rp := range proj.Roots {
		baseRoots[i] = rp.Ast
	}

	unreachable := make(map[string]struct{}, len(proj.AllPackages))
	for _, p := range proj.AllPackages {
		unreachable[p.Ast.PkgPath] = struct{}{}
	}
	for p := range packages.Postorder(baseRoots) {
		delete(unreachable, p.PkgPath)
	}
	return removePackages(slices.Collect(maps.Keys(unreachable)), proj)
}
