package writer

import (
	"path/filepath"

	"github.com/Grant-Nelson/Gozer/avail/coder"
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/modeler/remodel"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type Config struct {
	ErrGroup *faults.ErrGroup
}

type writerFactory struct {
	errGroup *faults.ErrGroup
}

func New(cfg *Config) remodel.RemodelFactory {
	return &writerFactory{
		errGroup: cfg.ErrGroup,
	}
}

func (f *writerFactory) StartPackage(pkg *project.Package) (bool, remodel.Remodeler, error) {
	w := &writer{
		errGroup: f.errGroup,
		pkg:      pkg,
	}
	return true, w, nil
}

type writer struct {
	errGroup *faults.ErrGroup
	pkg      *project.Package
	c        *coder.Coder
}

func (w *writer) getModuleFilePath(pkgPath string) string {
	// TODO: Need to properly create the file path to output.
	// The name probably should be part of the IR so that caching and writing
	// can always check the same path.
	return filepath.Join(pkgPath, `temp.ts`)
}

func (w *writer) Remodel() (bool, error) {
	// TODO: Add recover

	modulePath := w.getModuleFilePath(w.pkg.PkgPath())
	w.c = coder.New(modulePath, `  `, `//`)

	w.writePackageComment()
	w.writeImports()

	// TODO: Finish implementing

	return true, w.errGroup.FullOrNil()
}

func (w *writer) PackageDone() (bool, error) {
	// TODO: Add recover

	// TODO: Implement saving module file

	return true, nil
}

func (w *writer) writePackageComment() {
	if comment := w.pkg.Ir.Comment; len(comment) > 0 {
		w.c.WriteComment(comment)
		w.c.WriteLn()
	}
}

func (w *writer) writeImports() {
	imports := w.pkg.Ir.Imports
	for _, imp := range imports {
		w.writeImport(imp)
	}
	if len(imports) > 0 {
		w.c.WriteLn()
	}
}

func (w *writer) writeImport(imp *ir.ImportDef) {
	modulePath := w.getModuleFilePath(imp.PkgObj.Pkg().Path())
	name := imp.PkgObj.Name()
	w.c.WriteF(`import { %s } from '%s'`, name, modulePath)
	w.c.WriteLn()
}
