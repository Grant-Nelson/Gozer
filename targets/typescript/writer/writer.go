package writer

import (
	"github.com/Grant-Nelson/Gozer/avail/coder"
	"github.com/Grant-Nelson/Gozer/avail/faults"
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

func (w *writer) Remodel() (bool, error) {
	//ir := w.pkg.Ir
	//w.c.Write() ir.Comment

	// TODO: Implement

	return true, w.errGroup.FullOrNil()
}

func (w *writer) PackageDone() (bool, error) {

	// TODO: Implement saving module file

	return true, nil
}
