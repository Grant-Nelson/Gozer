package irc

import (
	"go/ast"
)

type Package struct {
	Funcs []*Func
}

// Func represents a function block defining a function as
// a collection of statement blocks.
// See [README.md]
type Func struct {
	// Ast is the original AST for this function.
	Ast *ast.FuncDecl

	// Blocks is the collection of statement blocks for this function.
	// The first block in this slice is the entry point for this function.
	Blocks []*Block
}

// Block represents a statement block defining a set of statements.
// Typically a statement block should have no flow control (mainly loops)
// except for exiting the block or for inlined flow (like a small if-statement).
// See [README.md]
type Block struct {
	// Index is the ones-based index of this block in the full set of blocks
	// for a package. 0 or less if not assigned an index yet.
	Index int

	// Ast are the original statements for the block.
	Ast []ast.Stmt

	// Prior blocks are blocks that can transition to this block.
	Prior []*Block

	// Follow blocks are blocks that this block can transition to.
	Follow []*Block

	// ExitsFunc indicates if this block can exit the function.
	ExitsFunc bool
}
