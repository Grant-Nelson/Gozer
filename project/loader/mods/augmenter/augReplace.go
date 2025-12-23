package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augReplace struct {
	fileSet *artifacts.FileSet
}

var _ mods.Modifier = (*augReplace)(nil)
var _ mods.LoadDoneExt = (*augReplace)(nil)

func (a *augReplace) Modify(f *artifacts.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augReplace) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
