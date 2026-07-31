package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// BlockRef is a reference for a block invocation.
// See [docs/Blocks.md] for more information.
type BlockRef struct {
	// Block is the block being referenced for invocation.
	// The reference is invalid when the block is null.
	Block *Block

	// Args are the arguments to pass onto the block when invoked.
	Args []ast.Expr
}

var _ Node = (*BlockRef)(nil)

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

func (b *BlockRef) Pos() token.Pos {
	if b == nil || len(b.Args) <= 0 || b.Args[0] == nil {
		return token.NoPos
	}
	return b.Args[0].Pos()
}
