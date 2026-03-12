package irc

import (
	"fmt"
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/assert"
)

// Package is the IRC for a whole package.
type Package struct {

	// Funcs is the collection of functions for this package.
	Funcs []*Func

	// TODO: Add types, vars, consts, etc
}

func (p *Package) String() string {
	return fmt.Sprintf("package{\n%s\n}", linesString(p.Funcs, `  `))
}

func (p *Package) NewFunc(astFunc *ast.FuncDecl) *Func {
	assert.NotNil(p)
	assert.NotNil(astFunc)

	fn := &Func{
		Ast:  astFunc,
		Name: astFunc.Name.Name,
	}
	p.Funcs = append(p.Funcs, fn)

	// Create initial block and populate it with current statements.
	block := fn.NewBlock()
	block.Hint = `initial`
	if astFunc.Body != nil {
		for _, s := range astFunc.Body.List {
			block.Body = append(block.Body, &BaseStmt{Stmt: s})
		}
	}
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
