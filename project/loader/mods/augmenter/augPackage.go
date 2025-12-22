package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augPackage struct {
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

func newPackage(build []string, basePath, testPkgPath string, fileSet *artifacts.FileSet) *augPackage {
	ap := &augPackage{
		del: &augDel{fileSet: fileSet},
		rep: &augReplace{fileSet: fileSet},
		ren: &augRename{fileSet: fileSet},
		add: &augAdd{fileSet: fileSet},

		build:       build,
		basePath:    basePath,
		testPkgPath: testPkgPath,
		fileSet:     fileSet,
	}
	ap.Group = mods.Group{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}

func (ap *augPackage) AddFile(filename string, src []byte, errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	ar := &augReader{augPackage: ap, errGroup: errGroup}
	ar.addFile(filename, src)
	return nil
}

func (ap *augPackage) LoadDone(errGroup *faults.Group) error {
	return ap.Group.LoadDone(errGroup)
}
