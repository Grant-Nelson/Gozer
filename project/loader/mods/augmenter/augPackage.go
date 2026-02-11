package augmenter

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

type augPackage struct {
	pkg *project.Package
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd
	mods.Group
}

var (
	_ mods.Modifier        = (*augPackage)(nil)
	_ mods.ModifyFileExt   = (*augPackage)(nil)
	_ mods.PackageStartExt = (*augPackage)(nil)
	_ mods.PackageDoneExt  = (*augPackage)(nil)
)

func newPackage(pkg *project.Package) *augPackage {
	ap := &augPackage{
		pkg: pkg,
		del: newDel(pkg),
		rep: newReplace(pkg),
		ren: newRename(pkg),
		add: newAdd(pkg),
	}
	ap.Group = mods.Group{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}

func (a *augPackage) ModName() string { return `Augmenter.Package` }

func (a *augPackage) AddFile(build []string, parser parser.Parser, filename string, src []byte, errGroup *faults.Group) (err error) {
	defer faults.Recover(&err) // TODO: Fix to use errGroup
	ar := newReader(a, build, parser, errGroup)
	ar.addFile(filename, src)
	return nil
}
