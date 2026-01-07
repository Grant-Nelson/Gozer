package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augPackage struct {
	pkg *artifacts.Package
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd
	mods.Group
}

var _ mods.Modifier = (*augPackage)(nil)
var _ mods.LoadDoneExt = (*augPackage)(nil)

func newPackage(pkg *artifacts.Package) *augPackage {
	ap := &augPackage{
		pkg: pkg,
		del: newDel(),
		rep: newReplace(),
		ren: newRename(),
		add: newAdd(pkg),
	}
	ap.Group = mods.Group{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}

func (a *augPackage) AddFile(build []string, filename string, src []byte, errGroup *faults.Group, fileParser artifacts.FileParser) (err error) {
	defer faults.Recover(&err)
	ar := newReader(a, errGroup, build, fileParser)
	ar.addFile(filename, src)
	return nil
}
