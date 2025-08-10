package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

type augReplace struct {
	fileSet *token.FileSet
}

func (a *augReplace) reset(fileSet *token.FileSet) {
	a.fileSet = fileSet
}

func (a *augReplace) Modify(fm *fileMod.FileMod, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augReplace) PackageDone(name, path string, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
