package augmenter

import (
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augReplace struct {
	fSet *token.FileSet
	pkg  *artifacts.Package
}

func newReplace(fSet *token.FileSet, pkg *artifacts.Package) *augReplace {
	return &augReplace{
		fSet: fSet,
		pkg:  pkg,
	}
}

var _ mods.Modifier = (*augReplace)(nil)
var _ mods.LoadDoneExt = (*augReplace)(nil)

func (a *augReplace) Modify(f *ast.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augReplace) LoadDone(errGroup *faults.Group) error {
	// TODO: Implement
	return nil
}
