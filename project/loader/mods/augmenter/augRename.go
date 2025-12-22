package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augRename struct {
	fileSet *artifacts.FileSet
}

func (a *augRename) Modify(f *artifacts.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
