package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augRename struct {
	fileSet *artifacts.FileSet
}

var _ mods.Modifier = (*augRename)(nil)
var _ mods.LoadDoneExt = (*augRename)(nil)

func (a *augRename) Modify(f *artifacts.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augRename) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
