package augmenter

import (
	"github.com/Grant-Nelson/Gozer/project/mods"
)

type Augmenter struct {
	mods.Group
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd
}

func New() *Augmenter {
	a := &Augmenter{
		del: &augDel{},
		rep: &augReplace{},
		ren: &augRename{},
		add: &augAdd{},
	}
	a.Group = mods.Group{a.del, a.rep, a.ren, a.add}
	return a
}

func (a *Augmenter) PackageStart(name, path string) error {
	// TODO: Populate the augmenter with data for this package.
	return a.Group.PackageStart(name, path)
}
