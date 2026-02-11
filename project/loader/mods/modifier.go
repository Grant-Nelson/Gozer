package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
)

// Modifier is a tool for modifying parts of a project during a load.
type Modifier interface {
	ModName() string
}

// ModifyFileExt extends a modifier to perform changes to the given file.
// If the modifier returns true, the modification will continue,
// otherwise all following modifiers will be skipped.
type ModifyFileExt interface {
	ModifyFile(f *ast.File, errGroup *faults.Group) (bool, error)
}

// PackageStartExt extends a modifier to indicate a package is starting.
// If the modifier returns true, the modification will continue,
// otherwise all following modifiers will be skipped.
type PackageStartExt interface {
	PackageStart(pkg *project.Package, errGroup *faults.Group) (bool, error)
}

// PackageDoneExt extends a modifier to indicate a package is done.
// If the modifier returns true, the modification will continue,
// otherwise all following modifiers will be skipped.
type PackageDoneExt interface {
	PackageDone(pkg *project.Package, errGroup *faults.Group) (bool, error)
}
