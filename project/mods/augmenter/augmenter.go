package augmenter

import (
	"github.com/Grant-Nelson/Gozer/project/fileMod"
	"github.com/Grant-Nelson/Gozer/project/mods"
)

type Augmenter struct {
	pkgPath string
	del     *augDel
	add     *augAdd
	rep     *augReplace
	ren     *augRename
}

func New() *Augmenter {
	return &Augmenter{
		del: &augDel{},
		add: &augAdd{},
		rep: &augReplace{},
		ren: &augRename{},
	}
}

func (a *Augmenter) Modify(fm *fileMod.FileMod) error {
	return mods.Modify(fm, a.del, a.add, a.rep, a.ren)
}

func (a *Augmenter) Finished() error {
	return mods.Finished(a.del, a.add, a.rep, a.ren)
}
