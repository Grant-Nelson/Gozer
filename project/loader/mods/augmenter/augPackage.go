package augmenter

import (
	"go/token"

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

func newPackage(fSet *token.FileSet, pkg *artifacts.Package) *augPackage {
	ap := &augPackage{
		pkg: pkg,
		del: newDel(fSet, pkg),
		rep: newReplace(fSet, pkg),
		ren: newRename(fSet, pkg),
		add: newAdd(fSet, pkg),
	}
	ap.Group = mods.Group{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}

func (a *augPackage) AddFile(build []string, fSet *token.FileSet, parser artifacts.Parser, filename string, src []byte, errGroup *faults.Group) (err error) {
	defer faults.Recover(&err)
	ar := newReader(a, build, parser, errGroup)
	ar.addFile(fSet, filename, src)
	return nil
}
