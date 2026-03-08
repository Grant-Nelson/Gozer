package remodel

import (
	"github.com/Grant-Nelson/Gozer/project"
	"github.com/Grant-Nelson/Gozer/project/modeler/irc"
)

type RemodelFactory interface {
	StartPackage(pkg *project.Package) (bool, Remodeler, error)
}

type Remodeler interface {
	PackageDone() (bool, error)
}

type RemodelFuncExt interface {
	RemodelFunc(f *irc.Func) (bool, error)
}
