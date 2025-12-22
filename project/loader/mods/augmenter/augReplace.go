package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augReplace struct {
	fileSet *artifacts.FileSet
}

func (a *augReplace) reset() {
	// TODO: Implement
}

func (a *augReplace) Modify(f *artifacts.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augReplace) PackageDone(pkg *artifacts.Package, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
