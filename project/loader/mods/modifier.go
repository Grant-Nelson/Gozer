package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

// Modifier performs a set changes to the given file.
// If the modifier returns true, the modification will continue, otherwise
// all following modifiers will be skipped.
type Modifier interface {
	Modify(f *ast.File, errGroup *faults.Group) (bool, error)
}

// LoadDoneExt extends a modifier to indicate the loading is done.
type LoadDoneExt interface {
	LoadDone(errGroup *faults.Group) error
}
