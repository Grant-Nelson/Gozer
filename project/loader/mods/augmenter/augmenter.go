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

	// ErrGroup is to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

type Augmenter struct {
	build    []string
	pathConv source.Converter
	parser   parser.Parser
	errGroup *faults.ErrGroup
}

var _ mods.ModFactory = (*Augmenter)(nil)

// Creates a new Modifier for augmenting Go files.
func New(cfg *Config) *Augmenter {
	return &Augmenter{
		build:    cfg.Build,
		pathConv: cfg.Converter,
		parser:   cfg.Parser,
		errGroup: cfg.ErrGroup,
	}
}

func (a *Augmenter) StartPackage(pkg *project.Package) (con bool, mod mods.Modifier, err error) {
	defer a.errGroup.Recover(&err)

	hasAug, augPath, augData, err := a.pathConv(pkg.PkgPath(), nil)
	if !hasAug {
		// Skip this package since there is no aug source.
		return true, nil, nil
	}

	augPkg := newPackage(pkg, a.errGroup)
	ar := newReader(augPkg, a.build, a.parser, a.errGroup)
	ar.readPackage(augPath, augData)
	return true, augPkg, nil
}
