package project

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
	"golang.org/x/tools/go/ssa"

	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/project/enums/buildState"
)

// Package describes a single package in a project.
type Package struct {
	// State indicates the status of the package during a build.
	State buildState.BuildState

	// Root indicates this package is a root package.
	Root bool

	// Imports is the packages that this package imports.
	Imports []*Package

	// Ast is the package in the Go abstract syntax tree (AST) form.
	//
	// This is either parsed from Go files and [Syntax] and [Types] is populated,
	// or this was loaded from the cache so only [Types] is populated.
	// This package is for the project's build, go version,
	// and augmented for the target languages.
	Ast *packages.Package

	// Saa is the package in the Go SSA IR form.
	Ssa *ssa.Package

	// Ir is the package's intermediate representation (IR) form.
	//
	// This is used to describe the code in a control flow graph (CFG) that
	// is used for optimization and translation to target language.
	Ir *ir.Package

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
	pkg := &Package{
		State: buildState.Listed,
		Ast:   basePkg,
	}
	for _, file := range pkg.Ast.Syntax {
		cleanFile(file)
	}
	return pkg
}

// cleanFile removes all the deprecated Object/Scope values
// to ensure that those values are never used.
func cleanFile(f *ast.File) {
	if f != nil {
		ast.Inspect(f, func(n ast.Node) bool {
			if id, ok := n.(*ast.Ident); ok {
				id.Obj = nil
			}
			return true
		})
		f.Scope = nil
		f.Unresolved = nil
	}
}

// PkgPath is the package path as used by the go/types package.
func (p *Package) PkgPath() string { return p.Ast.PkgPath }

// PkgName is the short name of the package.
// The package name is not guaranteed to be unique per program.
func (p *Package) PkgName() string { return p.Ast.Name }

// ForTest reports if this package is being built as part of a test.
func (p *Package) ForTest() bool { return len(p.Ast.ForTest) > 0 }

// IsTest reports if this package is being tested.
// This will be true for a X_test packages too.
func (p *Package) IsTest() bool {
	return p.IsXTest() || p.Ast.ForTest == p.PkgName()
}

// IsXTest reports if this package is a package specifically for testing a package,
// e.g., if `X` is the package, any file in the same path labelled `X_test`
// will comprise a package used for testing `X`.
func (p *Package) IsXTest() bool {
	return strings.HasSuffix(p.PkgName(), `_test`)
}

// FileSet is the set of files for positional information.
func (p *Package) FileSet() *token.FileSet { return p.Ast.Fset }

// Position converts a [Pos] location indicator into a Position value
// containing the file name, line number, column, and byte offset in file
// for the source file that the code labelled with the [Pos] came from.
func (p *Package) Position(pos token.Pos) token.Position {
	return p.FileSet().Position(pos)
}

// TODO: Need to be able to return composure and other
// information to help with optimization and blocking.
