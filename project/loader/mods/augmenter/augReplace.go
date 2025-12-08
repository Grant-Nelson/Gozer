package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/file"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type augReplace struct {
	fileSet *token.FileSet
}

func (a *augReplace) reset() {
	// TODO: Implement
}

func (a *augReplace) Modify(f *file.File, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augReplace) PackageDone(pkg *mods.Package, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
