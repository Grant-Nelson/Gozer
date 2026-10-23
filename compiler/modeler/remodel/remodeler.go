package remodel

import (
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type RemodelFactory interface {
	StartPackage(pkg *project.Package) (bool, Remodeler, error)
}

type Remodeler interface {
	Remodel() (bool, error)
}

type ProjectDoneExt interface {
	PackageDone() (bool, error)
}
