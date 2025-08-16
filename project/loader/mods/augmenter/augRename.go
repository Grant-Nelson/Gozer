package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type augRename struct {
	fileSet *token.FileSet
}

func (a *augRename) reset(fileSet *token.FileSet) {
	a.fileSet = fileSet
}

func (a *augRename) Modify(f *file.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) PackageDone(pkg *mods.Package, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) AddFile(f *file.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
