package augmenter

import (
	"go/ast"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

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

type Augmenter struct {
	build    []string
	pathConv PathConverter
	parser   parser.Parser
	curPkg   *augPackage
}

var (
	_ mods.Modifier        = (*Augmenter)(nil)
	_ mods.PackageStartExt = (*Augmenter)(nil)
	_ mods.PackageDoneExt  = (*Augmenter)(nil)
)

// Creates a new Modifier for augmenting Go files.
//
//   - The given build is the build constraints to load with.
//   - The pathConv is the conversion from the source paths to the augmentation files' paths.
//   - The parser is how files should be parsed and loaded.
//     If nil, the default file parser in the artifacts package.
func New(build []string, pathConv PathConverter, parser parser.Parser) *Augmenter {
	return &Augmenter{
		build:    build,
		pathConv: pathConv,
		parser:   parser,
		curPkg:   nil,
	}
}

func (a *Augmenter) Modify(f *ast.File, errGroup *faults.Group) (con bool, err error) {
	//defer faults.Recover(&err)// TODO: Connect these faults.Recover to errGroup
	if a.curPkg == nil {
		// no augmentation for this package.
		return true, nil
	}
	if con, err := a.curPkg.Modify(f, errGroup); err != nil || !con {
		return false, err
	}
	return true, nil
}

func (a *Augmenter) PackageStart(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	//defer faults.Recover(&err) // TODO: Connect these faults.Recover to errGroup
	// TODO: assert `a.curPkg == nil`

	hasAug, augPath := a.pathConv(pkg.PkgPath)
	if !hasAug {
		// store as nil to skip this package.
		return true, nil
	}

	a.curPkg = newPackage(pkg)
	ar := newReader(a.curPkg, a.build, a.parser, errGroup)
	ar.readPackage(augPath)
	return true, nil
}

func (a *Augmenter) PackageDone(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	//defer faults.Recover(&err)// TODO: Connect these faults.Recover to errGroup
	// TODO: assert `a.curPkg.pkg == pkg`

	if a.curPkg != nil {
		a.curPkg.PackageDone(pkg, errGroup)
		a.curPkg = nil
	}
	return true, nil
}
