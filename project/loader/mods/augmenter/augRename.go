package augmenter

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
	"github.com/Grant-Nelson/Gozer/project/loader/astMod"
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

func (a *augRename) PackageDone(pkg *astMod.PackageMod, errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}

func (a *augRename) AddFile(f *ast.File, errs *faults.Group) error {
	// TODO: Implement
	return nil
}
