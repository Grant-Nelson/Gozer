package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augReplace struct {
	fileSet *artifacts.FileSet
}

func (a *augReplace) Modify(f *artifacts.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augReplace) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
