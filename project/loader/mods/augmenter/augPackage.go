package augmenter

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/loader/mods"
	"github.com/Grant-Nelson/Gozer/project/loader/parser"
)

type augPackage struct {
	pkg      *project.Package
	errGroup *faults.ErrGroup

	del *augDel
	rep *augReplace
	ren *augRename
	add *augAdd
	mg  mods.ModGroup
}

var (
	_ mods.Modifier         = (*augPackage)(nil)
	_ mods.ModifyAstFileExt = (*augPackage)(nil)
)

func newPackage(pkg *project.Package, errGroup *faults.ErrGroup) *augPackage {
	ap := &augPackage{
		pkg:      pkg,
		errGroup: errGroup,

		del: newDel(pkg, errGroup),
		rep: newReplace(pkg, errGroup),
		ren: newRename(pkg, errGroup),
		add: newAdd(pkg, errGroup),
	}
	ap.mg = mods.ModGroup{ap.del, ap.rep, ap.ren, ap.add}
	return ap
}

func (a *augPackage) AddFile(build []string, parser parser.Parser, filename string, src []byte) (err error) {
	defer a.errGroup.Recover(&err)
	ar := newReader(a, build, parser, a.errGroup)
	ar.addFile(filename, src)
	return nil
}

func (a *augPackage) ModifyAstFile(f *ast.File) (con bool, err error) {
	defer a.errGroup.Recover(&err)
	if con, err := a.mg.ModifyAstFile(f); err != nil || !con {
		return false, err
	}
	return true, nil
}

func (a *augPackage) PackageDone() (con bool, err error) {
	defer a.errGroup.Recover(&err)
	return a.mg.PackageDone()
}
