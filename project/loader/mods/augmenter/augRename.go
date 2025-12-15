package augmenter

import (
	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augRename struct {
	fileSet *artifacts.FileSet
}

func (a *augRename) reset() {
	// TODO: Implement
}

func (a *augRename) Modify(f *artifacts.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) PackageDone(pkg *artifacts.Package, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
