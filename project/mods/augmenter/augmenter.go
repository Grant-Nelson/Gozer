package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/mods"
)

type Augmenter struct {
	mods.Group
	del     *augDel
	rep     *augReplace
	ren     *augRename
	add     *augAdd
	fileSet *token.FileSet
}

// TODO: Add build constraint information.
func New() *Augmenter {
	a := &Augmenter{
		del: &augDel{},
		rep: &augReplace{},
		ren: &augRename{},
		add: &augAdd{},
	}
	a.reset()
	a.Group = mods.Group{a.del, a.rep, a.ren, a.add}
	return a
}

func (a *Augmenter) reset() {
	a.fileSet = token.NewFileSet()
	a.del.reset(a.fileSet)
	a.rep.reset(a.fileSet)
	a.ren.reset(a.fileSet)
	a.add.reset(a.fileSet)
}

func (a *Augmenter) PackageStart(name, path string, errGroup *faults.Group) error {
	// TODO: Populate the augmenter with data for this package.
	return a.Group.PackageStart(name, path, errGroup)
}
