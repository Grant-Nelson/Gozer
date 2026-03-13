package remodel

import (
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/ir"
)

type RemodelFactory interface {
	StartPackage(pkg *project.Package) (bool, Remodeler, error)
}

type Remodeler interface {
	PackageDone() (bool, error)
}

type RemodelFuncExt interface {
	RemodelFunc(f *ir.Func) (bool, error)
}
