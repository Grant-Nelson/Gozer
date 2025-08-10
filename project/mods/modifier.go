package mods

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/internal/faults"
	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

// ErrFileModDone can be returned from a modifier to skip running
// any following modifiers for a group method call.
var ErrFileModDone = errors.New(`file modification is done`)

type (
	// Modifier performs a set changes to the given file.
	Modifier interface {
		Modify(fm *fileMod.FileMod, errGroup *faults.Group) error
	}

	// PackageStartExt extends a modifier to indicate the package has started loading.
	PackageStartExt interface {
		PackageStart(name, path string, errGroup *faults.Group) error
	}

	// PackageDoneExt extends a modifier to indicate the package is done loading.
	PackageDoneExt interface {
		PackageDone(name, path string, errGroup *faults.Group) error
	}

	// LoadDoneExt extends a modifier to indicate the loading is done.
	LoadDoneExt interface {
		LoadDone(errGroup *faults.Group) error
	}
)
