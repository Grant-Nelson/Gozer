package project

import (
	"go/token"
	"strings"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
	"github.com/Grant-Nelson/Gozer/project/interRep/irc"
	"golang.org/x/tools/go/packages"
)

// Package describes a single package in a project.
type Package struct {

	// State indicates the status of the package during a build.
	State buildState.BuildState

	// Root indicates this package is a root package.
	Root bool

	// Ast is the package in the Go abstract syntax tree (AST) form.
	//
	// This is either parsed from Go files and [Syntax] and [Types] is populated,
	// or this was loaded from the cache so only [Types] is populated.
	// This package is for the project's build, go version,
	// and augmented for the target languages.
	Ast *packages.Package

	// Irc is the package's intermediate representation code (IRC) form.
	//
	// This is used to describe the code in a control flow graph (CFG) that
	// is used for optimization and translation to target language.
	Irc *irc.Package

	// Depth is the depth of this node in the dependency tree where 0
	// means the package is a leave package with no dependencies and
	// the highest depth is a root. With multiple roots, some of the roots
	// may not have the highest depth value.
	//
	// This is used for determining which packages can be loaded in parallel.
	Depth int

	// TempTypeFile is set when the [types.Package] and related data was
	// serialized into a temporary file to be used as part of the cache.
	// This will be set to the path of that temporary file.
	TempTypeFile string
}

func newPackage(basePkg *packages.Package) *Package {
	return &Package{
		State: buildState.Listed,
		Ast:   basePkg,
	}
}

// PkgPath is the package path as used by the go/types package.
func (p *Package) PkgPath() string {
	return p.Ast.PkgPath
}

func (p *Package) IsTest() bool {
	return len(p.Ast.ForTest) > 0
}

func (p *Package) IsXTest() bool {
	return strings.HasSuffix(p.Ast.Name, `_test`)
}

// Position converts a [Pos] location indicator into a Position value
// containing the file name, line number, column, and byte offset in file
// for the source file that the code labelled with the [Pos] came from.
func (p *Package) Position(pos token.Pos) token.Position {
	return p.Ast.Fset.Position(pos)
}

// CollectErrors checks the types for any errors and
// adds all of them to the given error group.
func (p *Package) CollectErrors(errGroup *faults.ErrGroup) {
	errs := make([]error, len(p.Ast.Errors))
	for i, err := range p.Ast.Errors {
		errs[i] = err
	}
	errGroup.Add(errs...)
}
