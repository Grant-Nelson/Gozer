package project

import (
	"go/ast"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"

	"github.com/Grant-Nelson/Gozer/avail/faults"
	"github.com/Grant-Nelson/Gozer/compiler/ir"
	"github.com/Grant-Nelson/Gozer/compiler/project/enums/buildState"
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

func cleanFile(f *ast.File) {
	// Remove all the deprecated Object/Scope values.
	ast.Inspect(f, func(n ast.Node) bool {
		if id, ok := n.(*ast.Ident); ok {
			id.Obj = nil
		}
		return true
	})
	f.Scope = nil
}

// PkgPath is the package path as used by the go/types package.
func (p *Package) PkgPath() string {
	if p.Ir != nil {
		return p.Ir.Path
	}
	return p.Ast.PkgPath
}

// PkgName is the short name of the package.
// The package name is not guaranteed to be unique per program.
func (p *Package) PkgName() string {
	if p.Ir != nil {
		return p.Ir.Name
	}
	return p.Ast.Name
}

// IsTest reports if this package is being built as part of a test.
// This uses the AST so will always return false if the AST is nil.
func (p *Package) IsTest() bool {
	// TODO: Need to figure out a better way to do this so that `_test` files
	// are only included when the current package is the package being tested.
	return p.Ast != nil && len(p.Ast.ForTest) > 0
}

// IsXTest reports if this package is a package specifically for testing a package,
// e.g., if `X` is the package, any file in the same path labelled `X_test`
// will comprise a package used for testing `X`.
func (p *Package) IsXTest() bool {
	return strings.HasSuffix(p.PkgName(), `_test`)
}

func (p *Package) FileSet() *token.FileSet {
	if p.Ir != nil {
		return p.Ir.FileSet
	}
	return p.Ast.Fset
}

// Position converts a [Pos] location indicator into a Position value
// containing the file name, line number, column, and byte offset in file
// for the source file that the code labelled with the [Pos] came from.
func (p *Package) Position(pos token.Pos) token.Position {
	return p.FileSet().Position(pos)
}

// CollectErrors checks the types for any errors in the AST package and
// adds all of them to the given error group.
func (p *Package) CollectErrors(errGroup *faults.ErrGroup) {
	if p.Ast != nil {
		errs := make([]error, len(p.Ast.Errors))
		for i, err := range p.Ast.Errors {
			errs[i] = err
		}
		errGroup.Add(errs...)
	}
}
