package ir

import (
	"go/token"
	"go/types"
)

// Package is the IR for a whole package.
type Package struct {
	// FileSet is the information about the source files from the AST.
	FileSet *token.FileSet

	// Name is the short name for this package.
	Name string

	// Path is the path that would be used to import this package.
	Path string

	// Sizes are the sizes for the architecture.
	// For example TS and JS will have 32-bit `int`, `uint`, and `uintptr`.
	Sizes types.Sizes

	// Imports are the packages that are imported by this package.
	Imports map[string]*ImportDecl

	// Types is the collection of types for this package.
	Types []*TypeDecl

	// Consts is the collection of const values for this package.
	Consts []*ConstDecl

	// Vars is the collection of variables for this package.
	Vars []*VarDecl

	// Funcs is the collection of functions for this package.
	Funcs []*FuncDecl
}

var _ Parent = (*Package)(nil)

func (p *Package) String() string {
	result := `package{` + nlStr
	if len(p.Name) > 0 {
		result += indentStr + `Name: ` + p.Name + nlStr
	}
	if len(p.Path) > 0 && p.Path != `command-line-arguments` {
		result += indentStr + `Path: ` + p.Path + nlStr
	}

	// TODO: Add Imports
	if len(p.Types) > 0 {
		result += indentStr + `Types:` + indentInner(bodyString(p.Types)) + nlStr
	}
	if len(p.Consts) > 0 {
		result += indentStr + `Consts:` + indentInner(bodyString(p.Consts)) + nlStr
	}
	if len(p.Vars) > 0 {
		result += indentStr + `Vars:` + indentInner(bodyString(p.Vars)) + nlStr
	}
	if len(p.Funcs) > 0 {
		result += indentStr + `Funcs:` + indentInner(bodyString(p.Funcs)) + nlStr
	}
	return result + `}`
}

func (p *Package) Children(yield func(Node) bool) {
	_ = YieldSlice(p.Types, yield) && YieldSlice(p.Consts, yield) &&
		YieldSlice(p.Vars, yield) && YieldSlice(p.Funcs, yield)
}
