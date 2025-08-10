package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type augRename struct {
	fileSet *token.FileSet
}

func (a *augRename) reset(fileSet *token.FileSet) {
	a.fileSet = fileSet
}

func (a *augRename) Modify(fm *fileMod.FileMod, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) PackageDone(name, path string, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
