package remodel

import (
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/project"
)

type RemodelFactory interface {
	StartPackage(pkg *project.Package) (bool, Remodeler, error)
}

type Remodeler interface {
	PackageDone() (bool, error)
}

// TODO: Add Remodel Type, Var, and Const

type RemodelFuncExt interface {
	RemodelFunc(f *ir.Func) (bool, error)
}
