package augmenter

import (
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

type augPackage struct {
	mods.Group
	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd
	pkg *artifacts.Package
}

var _ mods.Modifier = (*augPackage)(nil)
var _ mods.LoadDoneExt = (*augPackage)(nil)

func newPackage(pkg *artifacts.Package) *augPackage {
	fs := pkg.TempFileSet()
	ap := &augPackage{
		del: &augDel{fileSet: fs},
		rep: &augReplace{fileSet: fs},
		ren: &augRename{fileSet: fs},
		add: &augAdd{fileSet: fs},
		pkg: pkg,
	}
	ap.Group = mods.Group{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}
