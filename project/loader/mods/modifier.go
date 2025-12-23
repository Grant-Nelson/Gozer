package mods

import (
	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/loader/mods/artifacts"
)

// Modifier performs a set changes to the given file.
// If the modifier returns true, the modification will continue, otherwise
// all following modifiers will be skipped.
type Modifier interface {
	Modify(f *artifacts.File, errGroup *faults.Group) (bool, error)
}

// LoadDoneExt extends a modifier to indicate the loading is done.
type LoadDoneExt interface {
	LoadDone(errGroup *faults.Group) error
}
