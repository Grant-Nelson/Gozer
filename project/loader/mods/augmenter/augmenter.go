package augmenter

import (
	"go/ast"
	"go/token"
	"maps"
	"slices"
	"sort"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type Augmenter struct {
	build    []string
	fSet     *token.FileSet
	pathConv PathConverter
	parser   artifacts.Parser
	packages map[string]*augPackage
}

// PathConverter takes the given import path for a package and returns
// the paths to where the augmentation files are found.
// Return false to skip this path.
//
// TODO: Maybe make this more full featured so that it can return `[]byte` data
// for a file instead of requiring it to be loaded from a file.
type PathConverter func(string) (bool, string)

// PathRebase performs a simple string-wise match of a prefix path (oldBase)
// and replaces it with a new prefix path (newBase).
// If the oldBase id not prefixed then it will return false.
// If oldBase is empty, this will simply always concatenate the newBase.
func PathRebase(oldBase, newBase string) PathConverter {
	if len(oldBase) <= 0 {
		return func(s string) (bool, string) {
			return true, newBase + s
		}
	}
	return func(s string) (bool, string) {
		if suffix, ok := strings.CutPrefix(s, oldBase); ok {
			return true, newBase + suffix
		}
		return false, ``
	}
}

var _ mods.Modifier = (*Augmenter)(nil)
var _ mods.LoadDoneExt = (*Augmenter)(nil)

// Creates a new Modifier for augmenting Go files.
//
//   - The given build is the build constraints to load with.
//   - The pathConv is the conversion from the source paths to the augmentation files' paths.
//   - The parser is how files should be parsed and loaded.
//     If nil, the default file parser in the artifacts package.
func New(build []string, fSet *token.FileSet, pathConv PathConverter, parser artifacts.Parser) *Augmenter {
	return &Augmenter{
		build:    build,
		fSet:     fSet,
		pathConv: pathConv,
		parser:   parser,
		packages: map[string]*augPackage{},
	}
}

func (a *Augmenter) Modify(f *ast.File, errGroup *faults.Group) (con bool, err error) {
	defer faults.Recover(&err)
	pkg := artifacts.PackageForFile(a.fSet, f)
	ap, exists := a.getPackage(pkg.Key())
	if !exists {
		ap = a.addPackage(pkg, errGroup)
	}
	if ap == nil {
		return true, nil
	}
	if con, err := ap.Modify(f, errGroup); err != nil || !con {
		return false, err
	}
	return true, nil
}

func (a *Augmenter) LoadDone(errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	keys := slices.Collect(maps.Keys(a.packages))
	sort.Strings(keys)
	for _, key := range keys {
		if ap := a.packages[key]; ap != nil {
			ap.LoadDone(errGroup)
		}
	}
	return nil
}

func (a *Augmenter) getPackage(pkgKey string) (*augPackage, bool) {
	ap, exists := a.packages[pkgKey]
	return ap, exists
}

func (a *Augmenter) addPackage(pkg *artifacts.Package, errGroup *faults.Group) *augPackage {
	key := pkg.Key()

	hasAug, augPath := a.pathConv(pkg.Path())
	if !hasAug {
		// store as nil to prevent trying again
		a.packages[key] = nil
		return nil
	}

	ap := newPackage(a.fSet, pkg)
	a.packages[key] = ap

	ar := newReader(ap, a.build, a.parser, errGroup)
	ar.readPackage(a.fSet, augPath)
	return ap
}
