package augmenter

import (
	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type Augmenter struct {
	mods.Group
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd

	build       []string
	basePath    string
	testPkgPath string
	fileSet     *artifacts.FileSet
}

func New(build []string, basePath, testPkgPath string, fileSet *artifacts.FileSet) *Augmenter {
	a := &Augmenter{
		del: &augDel{fileSet: fileSet},
		rep: &augReplace{fileSet: fileSet},
		ren: &augRename{fileSet: fileSet},
		add: &augAdd{fileSet: fileSet},

		build:       build,
		basePath:    basePath,
		testPkgPath: testPkgPath,
		fileSet:     fileSet,
	}
	a.Group = mods.Group{a.del, a.rep, a.ren, a.add}
	return a
}

func (a *Augmenter) PackageStart(pkg *artifacts.Package, errGroup *faults.Group) error {
	a.reset()
	if err := a.AddPackage(pkg.Path, errGroup); err != nil {
		return err
	}
	return a.Group.PackageStart(pkg, errGroup)
}

func (a *Augmenter) reset() {
	a.del.reset()
	a.rep.reset()
	a.ren.reset()
	a.add.reset()
}
