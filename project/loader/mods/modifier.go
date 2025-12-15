package mods

import (
	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

// Modifier performs a set changes to the given file.
type Modifier interface {
	Modify(f *artifacts.File, errGroup *faults.Group) error
}

// PackageStartExt extends a modifier to indicate the package has started loading.
type PackageStartExt interface {
	PackageStart(pkg *artifacts.Package, errGroup *faults.Group) error
}

// PackageDoneExt extends a modifier to indicate the package is done loading.
type PackageDoneExt interface {
	PackageDone(pkg *artifacts.Package, errGroup *faults.Group) error
}

// LoadDoneExt extends a modifier to indicate the loading is done.
type LoadDoneExt interface {
	LoadDone(errGroup *faults.Group) error
}
