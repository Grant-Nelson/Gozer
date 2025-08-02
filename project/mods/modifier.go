package mods

import (
	"errors"

	"github.com/Grant-Nelson/Gozer/project/fileMod"
)

// ErrFileModDone can be returned from a modifier to skip running
// any following modifiers for a group method call.
var ErrFileModDone = errors.New(`file modification is done`)

type (
	// Modifier performs a set changes to the given file.
	Modifier interface {
		Modify(fm *fileMod.FileMod) error
	}

	// PackageStartExt extends a modifier to indicate the package has started loading.
	PackageStartExt interface {
		PackageStart(name, path string) error
	}

	// PackageDoneExt extends a modifier to indicate the package is done loading.
	PackageDoneExt interface {
		PackageDone(name, path string) error
	}

	// LoadDoneExt extends a modifier to indicate the loading is done.
	LoadDoneExt interface {
		LoadDone() error
	}
)
