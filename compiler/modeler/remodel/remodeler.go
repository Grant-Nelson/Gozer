package remodel

import "github.com/Grant-Nelson/Gozer/compiler/project"

// RemodelFactory is a tool for creating remodelers to remodeling parts of a
// project's packages during a load.
//
// The factory needs to be able to be called concurrently
// so that the constructed Remodelers needs to be able to run asynchronously.
type RemodelFactory interface {

	// StartPackage constructs a new remodeler for the given package.
	//
	// If this returns true, the remodel of this package will continue,
	// otherwise the remodel of this package will skip starting the
	// package in all following remodeler factories.
	//
	// This may return a remodeler or may return nil if no other remodel
	// is needed from this remodel.
	StartPackage(pkg *project.Package) (bool, Remodeler, error)
}

// Remodeler will perform changes to the given package.
type Remodeler interface {

	// Remodel is called per package being remodeled.
	//
	// If this returns true, the remodeler of this package will continue,
	// otherwise all following remodelers will be skipped for this file only.
	Remodel() (bool, error)
}

// ProjectDoneExt extends a [Modifier] to perform a remodel on a package.
type ProjectDoneExt interface {

	// PackageDone is called when a package finishing its remodeling.
	//
	// This may not always be called based on the results from other remodelers.
	//
	// If this returns true, the following remodelers will continue to be called,
	// otherwise all following remodelers will be skipped.
	PackageDone() (bool, error)
}
