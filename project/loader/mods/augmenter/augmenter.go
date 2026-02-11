package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

type Augmenter struct {
	build    []string
	pathConv parser.SourceConverter
	parser   parser.Parser
	curPkg   *augPackage
}

var (
	_ mods.Modifier        = (*Augmenter)(nil)
	_ mods.ModifyFileExt   = (*Augmenter)(nil)
	_ mods.PackageStartExt = (*Augmenter)(nil)
	_ mods.PackageDoneExt  = (*Augmenter)(nil)
)

// Creates a new Modifier for augmenting Go files.
//
//   - The given build is the build constraints to load with.
//   - The pathConv is the conversion from the source paths to the augmentation files' paths.
//   - The parser is how files should be parsed and loaded.
//     If nil, the default file parser in the artifacts package.
func New(build []string, pathConv parser.SourceConverter, parser parser.Parser) *Augmenter {
	return &Augmenter{
		build:    build,
		pathConv: pathConv,
		parser:   parser,
		curPkg:   nil,
	}
}

func (a *Augmenter) ModName() string { return `Augmenter` }

func (a *Augmenter) ModifyFile(f *ast.File, errGroup *faults.Group) (con bool, err error) {
	defer faults.Recover(&err) // TODO: Connect these faults.Recover to errGroup
	if a.curPkg == nil {
		// no augmentation for this package.
		return true, nil
	}
	if con, err := a.curPkg.ModifyFile(f, errGroup); err != nil || !con {
		return false, err
	}
	return true, nil
}

func (a *Augmenter) PackageStart(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	defer faults.Recover(&err) // TODO: Connect these faults.Recover to errGroup
	// TODO: assert `a.curPkg == nil`

	hasAug, augPath, augData, err := a.pathConv(pkg.PkgPath, nil)
	if !hasAug {
		// store as nil to skip this package.
		return true, nil
	}

	a.curPkg = newPackage(pkg)
	ar := newReader(a.curPkg, a.build, a.parser, errGroup)
	ar.readPackage(augPath, augData)
	return true, nil
}

func (a *Augmenter) PackageDone(pkg *project.Package, errGroup *faults.Group) (con bool, err error) {
	defer faults.Recover(&err) // TODO: Connect these faults.Recover to errGroup
	// TODO: assert `a.curPkg.pkg == pkg`

	if a.curPkg != nil {
		a.curPkg.PackageDone(pkg, errGroup)
		a.curPkg = nil
	}
	return true, nil
}
