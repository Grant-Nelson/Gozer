package astMod

import (
	"go/token"

	"github.com/Grant-Nelson/Gozer/internal/faults"
)

type PackageMod struct {
	name        string
	path        string
	errGroup    *faults.Group
	tempFileSet *token.FileSet
	// TODO: Add indication of if this is a test package.
}

// NewPackage creates a new temporary package used when modifying files' AST.
func NewPackage(name, path string, errGroup *faults.Group, tempFileSet *token.FileSet) *PackageMod {
	return &PackageMod{
		name:        name,
		path:        path,
		errGroup:    errGroup,
		tempFileSet: tempFileSet,
	}
}

// Name is the name of this package.
func (pm *PackageMod) Name() string { return pm.name }

// Path is the import path of this package.
func (pm *PackageMod) Path() string { return pm.path }

// ErrorGroup is the error group used while loading
// to report several errors at a time.
func (pm *PackageMod) ErrorGroup() *faults.Group { return pm.errGroup }

// FileSet is used to set the tracing for the file.
// This is a temporary file set specific to storing this file
// and additional files while loading it.
func (pm *PackageMod) FileSet() *token.FileSet { return pm.tempFileSet }
