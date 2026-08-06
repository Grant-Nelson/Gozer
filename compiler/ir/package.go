package ir

import (
	"fmt"
	"go/token"
	"go/types"
)

// Package is the IR for a whole package.
type Package struct {
	// Info is the type information from the AST.
	Info *types.Info // TODO: REMOVE

	// FileSet is the information about the source files from the AST.
	FileSet *token.FileSet // TODO: REMOVE

	// Funcs is the collection of functions for this package.
	Funcs []*Func

	// TODO: Add types, vars, consts, etc
}

func (p *Package) String() string {
	return fmt.Sprintf("package{\n%s\n}", linesString(p.Funcs))
}

// FindFunc finds a function with the given name or nil if not found.
func (p *Package) FindFunc(name string) *Func {
	for _, fn := range p.Funcs {
		if fn.Name == name {
			return fn
		}
	}
	return nil
}
