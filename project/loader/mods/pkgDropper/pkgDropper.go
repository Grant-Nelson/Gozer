package pkgDropper

import (
	"fmt"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type Config struct {
	// ErrGroup is used to collect multiple errors.
	ErrGroup *faults.ErrGroup
}

type PkgDropper struct {
	errGroup *faults.ErrGroup
}

var _ mods.ModFactory = (*PkgDropper)(nil)

func New(cfg *Config) *PkgDropper {
	return &PkgDropper{
		errGroup: cfg.ErrGroup,
	}
}

func (fc *PkgDropper) StartPackage(pkg *project.Package) (bool, mods.Modifier, error) {

	// TODO: Finish implementing or remove, or rename to expunge specific packages.

	fmt.Printf("==[ %s ]========================\n", pkg.PkgPath())
	for i, file := range pkg.Ast.GoFiles {
		fmt.Printf(">> %d. %s\n", i, file)
	}
	return true, nil, nil
}
