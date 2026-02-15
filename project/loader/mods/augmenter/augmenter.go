package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/avail/source"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

type Config struct {
	// Build is the build constraints to load with.
	Build []string

	// Converter is the conversion from the source paths to the
	// augmentation files' paths.
	Converter source.Converter

	// Parser is the method for parsing Go files.
	// If nil, the default file parser.
	Parser parser.Parser
}

type Augmenter struct {
	build    []string
	pathConv source.Converter
	parser   parser.Parser
}

var _ mods.ModFactory = (*Augmenter)(nil)

// Creates a new Modifier for augmenting Go files.
func New(cfg *Config) *Augmenter {
	return &Augmenter{
		build:    cfg.Build,
		pathConv: cfg.Converter,
		parser:   cfg.Parser,
	}
}

func (a *Augmenter) StartPackage(pkg *project.Package, errGroup *faults.Group) (con bool, mod mods.Modifier, err error) {
	defer faults.Recover(&err) // TODO: Connect these faults.Recover to errGroup
	// TODO: assert `a.curPkg == nil`

	hasAug, augPath, augData, err := a.pathConv(pkg.PkgPath(), nil)
	if !hasAug {
		// Skip this package since there is no aug source.
		return true, nil, nil
	}

	augPkg := newPackage(pkg, errGroup)
	ar := newReader(augPkg, a.build, a.parser, errGroup)
	ar.readPackage(augPath, augData)
	return true, augPkg, nil
}
