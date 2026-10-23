package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/loader/mods"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type augRename struct {
	pkg      *project.Package
	errGroup *faults.ErrGroup
}

func newRename(pkg *project.Package, errGroup *faults.ErrGroup) *augRename {
	return &augRename{
		pkg:      pkg,
		errGroup: errGroup,
	}
}

var (
	_ mods.Modifier         = (*augRename)(nil)
	_ mods.ModifyAstFileExt = (*augRename)(nil)
)

func (a *augRename) ModifyAstFile(f *ast.File) (bool, error) {
	// TODO: Implement
	return true, nil
}

func (a *augRename) PackageDone() (bool, error) {
	// TODO: Implement
	return true, nil
}
