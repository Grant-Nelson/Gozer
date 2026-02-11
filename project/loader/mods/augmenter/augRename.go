package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type augRename struct {
	pkg *project.Package
}

func newRename(pkg *project.Package) *augRename {
	return &augRename{
		pkg: pkg,
	}
}

var (
	_ mods.Modifier       = (*augRename)(nil)
	_ mods.ModifyFileExt  = (*augRename)(nil)
	_ mods.PackageDoneExt = (*augRename)(nil)
)

func (a *augRename) ModName() string { return `Augmenter.Rename` }

func (a *augRename) ModifyFile(f *ast.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augRename) PackageDone(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}
