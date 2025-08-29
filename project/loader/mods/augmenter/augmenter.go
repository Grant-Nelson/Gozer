package augmenter

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
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
	fileSet     *token.FileSet
}

func New(build []string, basePath, testPkgPath string) *Augmenter {
	a := &Augmenter{
		del: &augDel{},
		rep: &augReplace{},
		ren: &augRename{},
		add: &augAdd{},

		build:       build,
		testPkgPath: testPkgPath,
		basePath:    basePath,
	}
	a.Group = mods.Group{a.del, a.rep, a.ren, a.add}
	return a
}

func (a *Augmenter) PackageStart(pkg *mods.Package, errGroup *faults.Group) error {
	a.reset()
	if err := a.AddPackage(pkg.Path, errGroup); err != nil {
		return err
	}
	return a.Group.PackageStart(pkg, errGroup)
}

func (a *Augmenter) reset() {
	a.fileSet = token.NewFileSet()
	a.del.reset(a.fileSet)
	a.rep.reset(a.fileSet)
	a.ren.reset(a.fileSet)
	a.add.reset(a.fileSet)
}
