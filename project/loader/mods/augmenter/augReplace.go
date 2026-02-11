package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
)

type augReplace struct {
	pkg *project.Package
}

func newReplace(pkg *project.Package) *augReplace {
	return &augReplace{
		pkg: pkg,
	}
}

var (
	_ mods.Modifier       = (*augReplace)(nil)
	_ mods.ModifyFileExt  = (*augReplace)(nil)
	_ mods.PackageDoneExt = (*augReplace)(nil)
)

func (a *augReplace) ModName() string { return `Augmenter.Replace` }

func (a *augReplace) ModifyFile(f *ast.File, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augReplace) PackageDone(pkg *project.Package, errGroup *faults.Group) (bool, error) {
	// TODO: Implement
	return true, nil
}
