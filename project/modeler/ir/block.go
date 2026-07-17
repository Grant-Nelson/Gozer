package ir

import (
	"fmt"
	"go/ast"
	"go/token"
	"strings"
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

		// Hint is a string to help debug block creation.
		Hint string

		// Body is the list of statements for this block.
		// A block is invalid when there are no statements.
		Body []Stmt

		// Prior blocks are blocks that can transition to this block.
		Prior []*Block

		// Follow blocks are blocks that this block can transition to.
		Follow []*Block

		// Params are the parameters that are passed into this block
		// when it is called and are available inside the block.
		//
		// This will not contain any parameters needed for the closure
		// including any receiver and any type dictionaries.
		Params []*Param
	}

	// BlockRef is a reference for a block invocation.
	// See [docs/Blocks.md] for more information.
	BlockRef struct {
		// Block is the block being referenced for invocation.
		// The reference is invalid when the block is null.
		Block *Block

		// Args are the arguments to pass onto the block when invoked.
		Args []ast.Expr
	}
)

var _ Node = (*Block)(nil)

func (b *Block) String() string {
	var params, hint, tail string
	if len(b.Params) > 0 {
		parts := make([]string, len(b.Params))
		for i, p := range b.Params {
			parts[i] = p.String()
		}
		params = strings.Join(parts, `, `)
	}
	if len(b.Hint) > 0 {
		hint = `<` + b.Hint + `> `
	}
	if len(b.Body) > 0 {
		tail = fmt.Sprintf("\n%s\n", linesString(b.Body))
	}
	return fmt.Sprintf("block %d (%s)%s{%s}", b.Index, params, hint, tail)
}

func (b *Block) Pos() token.Pos {
	if b == nil || len(b.Body) <= 0 || b.Body[0] == nil {
		return token.NoPos
	}
	return b.Body[0].Pos()
}

// LastStmt gets the last statement in the block or nil if empty.
func (b *Block) LastStmt() Stmt {
	if max := len(b.Body) - 1; max >= 0 {
		return b.Body[max]
	}
	return nil
}

func (r *BlockRef) String() string {
	var tail string
	if len(r.Args) > 0 {
		tail = fmt.Sprintf(`, [%s]`, csvString(r.Args))
	}
	if r.Block == nil {
		return fmt.Sprintf(`block <nil>%s`, tail)
	}
	return fmt.Sprintf(`block %d%s`, r.Block.Index, tail)
}
