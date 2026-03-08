package irc

import (
	"fmt"
	"go/ast"
	"strings"
)

type (
	// Package is the IRC for a whole package.
	Package struct {
		Funcs []*Func
	}

	// Func represents a function block defining a function as
	// a collection of statement blocks.
	// See [README.md]
	Func struct {
		// Ast is the AST for this function.
		//
		// Remodelers that modify the AST inside of a [BaseStmt] or similar
		// may cause changes to this data as well.
		Ast *ast.FuncDecl

		// Name of the function.
		Name string

		// Blocks is the collection of statement blocks for this function.
		// The first block in this slice is the entry point for this function.
		Blocks []*Block
	}
)

func (p *Package) String() string {
	return fmt.Sprintf("package{\n%s\n}", linesString(p.Funcs, `  `))
}

// NewBlock creates a new empty block and adds it to this function.
func (fn *Func) NewBlock() *Block {
	b := &Block{Index: len(fn.Blocks)}
	fn.Blocks = append(fn.Blocks, b)
	return b
}

func (fn *Func) String() string {
	return fmt.Sprintf("func %s {\n%s\n}", fn.Name, linesString(fn.Blocks, `  `))
}

func csvString[E any, S []E](s S) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = fmt.Sprintf(`%v`, elem)
	}
	return strings.Join(elems, `, `)
}

func linesString[E any, S []E](s S, indent string) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := fmt.Sprintf(`%s%v`, indent, elem)
		elems[i] = strings.ReplaceAll(eStr, "\n", "\n"+indent)
	}
	return strings.Join(elems, "\n")
}
