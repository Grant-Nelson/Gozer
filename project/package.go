package project

import (
	"go/token"
	"strings"

	"github.com/Grant-Nelson/Gozer/project/enums/buildState"
	"golang.org/x/tools/go/packages"
)

// Package describes a single package in a project.
type Package struct {

	// State indicates the status of the package during a build.
	State buildState.BuildState

	// Ast is the package in the Go AST form.
	//
	// This is either parsed from Go files and [Syntax] and [Types] is populated,
	// or this was loaded from the cache so only [Types] is populated.
	// This package is for the project's build, go version,
	// and augmented for the target languages.
	Ast *packages.Package

	// TempTypeFile is set when the [types.Package] and related data was
	// serialized into a temporary file to be used as part of the cache.
	// This will be set to the path of that temporary file.
	TempTypeFile string
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
