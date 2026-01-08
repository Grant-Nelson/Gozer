package augmenter

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augRename struct {
	fSet *token.FileSet
	pkg  *artifacts.Package
}

func newRename(fSet *token.FileSet, pkg *artifacts.Package) *augRename {
	return &augRename{
		fSet: fSet,
		pkg:  pkg,
	}
}

var _ mods.Modifier = (*augRename)(nil)
var _ mods.LoadDoneExt = (*augRename)(nil)

func (a *augRename) Modify(f *ast.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augRename) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
