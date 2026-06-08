package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/avail/assert"
)

// Package is the IR for a whole package.
type Package struct {
	// Info is the type information from the AST.
	Info *types.Info

	// FileSet is the information about the source files from the AST.
	FileSet *token.FileSet

	// Funcs is the collection of functions for this package.
	Funcs []*Func

	// TODO: Add types, vars, consts, etc
}

func (p *Package) String() string {
	return fmt.Sprintf("package{\n%s\n}", linesString(p.Funcs))
}

func (p *Package) NewFunc(astFunc *ast.FuncDecl) *Func {
	assert.NotNil(p)
	assert.NotNil(astFunc)

	fn := &Func{
		Ast:  astFunc,
		Name: astFunc.Name.Name,
	}
	p.Funcs = append(p.Funcs, fn)

	// Create initial block, populate it with current statements, add params.
	conv := &Converter{Info: p.Info}
	body := conv.ExpandStmt(conv.FromBlockStmt(astFunc.Body))
	params := conv.ExpandParams(astFunc.Type)
	fn.NewBlock(`initial`, body, params)
	return fn
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
