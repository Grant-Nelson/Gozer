package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type augReplace struct {
	pkg      *project.Package
	errGroup *faults.ErrGroup
}

func newReplace(pkg *project.Package, errGroup *faults.ErrGroup) *augReplace {
	return &augReplace{
		pkg:      pkg,
		errGroup: errGroup,
	}
}

var (
	_ mods.Modifier         = (*augReplace)(nil)
	_ mods.ModifyAstFileExt = (*augReplace)(nil)
)

func (a *augReplace) ModifyAstFile(f *ast.File) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augReplace) PackageDone() (bool, error) {
	// TODO: Implement
	return true, nil
}
