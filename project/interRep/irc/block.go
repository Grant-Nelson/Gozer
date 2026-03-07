package irc

import (
	"fmt"
	"go/ast"
)

type (
	// Block represents a statement block defining a set of statements.
	// Typically a statement block should have no flow control (mainly loops)
	// except for exiting the block or for inlined flow (like a small if-statement).
	// See [README.md]
	Block struct {
		// Index is the ones-based index of this block in the full set of blocks
		// for a package. 0 or less if not assigned an index yet.
		Index int

		// Body is the list of statements for this block.
		Body []Stmt

		// Prior blocks are blocks that can transition to this block.
		Prior []*Block

		// Follow blocks are blocks that this block can transition to.
		Follow []*Block

		// ExitsFunc indicates if this block can exit the function.
		ExitsFunc bool
	}

	// BlockRef is a reference for a block invocation.
	// See [docs/Blocks.md] for more information.
	BlockRef struct {
		Block *Block     // the block being referenced for invocation
		Args  []ast.Expr // the arguments to pass onto the block when invoked
	}
)

func (b *Block) String() string {
	return fmt.Sprintf("block %d {\n%s\n}", b.Index, linesString(b.Body, `  `))
}

// LastStmt gets the last statement in the block or nil if empty.
func (b *Block) LastStmt() Stmt {
	max := len(b.Body) - 1
	if max >= 0 {
		return b.Body[max]
	}
	return nil
}

func (r *BlockRef) String() string {
	return fmt.Sprintf(`block %d, [%s]`, r.Block.Index, csvString(r.Args))
}
