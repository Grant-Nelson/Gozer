package mods

import (
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project"
)

// ModFactory is a tool for creating modifiers to modifying parts of a
// project's packages during a load.
//
// The factory needs to be able to be called concurrently
// so that the constructed Modifiers needs to be able to run asynchronously.
type ModFactory interface {

	// StartPackage constructs a new modifier for the given package.
	//
	// If this returns true, the modification of this package will continue,
	// otherwise the modification of this package will skip starting the
	// package in all following modifier factories.
	//
	// This may return a modifier or may return nil if no other modification
	// is needed from this modification.
	StartPackage(pkg *project.Package, errGroup *faults.Group) (bool, Modifier, error)
}

// Modifier will perform a modification on a package.
type Modifier interface {

	// PackageDone is called when a package finishing its modification.
	//
	// This may not always be called based on the results from other modifiers.
	//
	// If this returns true, the following modifiers will continue to be called,
	// otherwise all following modifiers will be skipped.
	PackageDone() (bool, error)
}

// ModifyAstFileExt extends a [Modifier] to perform changes to the given AST file.
type ModifyAstFileExt interface {

	// ModifyAstFile is called per AST file being modified.
	//
	// If this returns true, the modification of this file will continue,
	// otherwise all following modifiers will be skipped for this file only.
	ModifyAstFile(f *ast.File) (bool, error)
}
