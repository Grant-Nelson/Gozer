package ir

import (
	"fmt"
	"go/token"
)

// BlockRef is a reference for a block invocation.
// See [docs/Blocks.md] for more information.
type BlockRef struct {
	// Block is the block being referenced for invocation.
	// The reference is invalid when the block is null.
	Block *Block

	// Args are the arguments to pass onto the block when invoked.
	Args []Expr
}

var (
	_ Node   = (*BlockRef)(nil)
	_ Parent = (*BlockRef)(nil)
)

func (n *BlockRef) String() string {
	var tail string
	if len(n.Args) > 0 {
		tail = fmt.Sprintf(`, [%s]`, csvString(n.Args))
	}
	if n.Block == nil {
		return fmt.Sprintf(`block <nil>%s`, tail)
	}
	return fmt.Sprintf(`block %d%s`, n.Block.Index, tail)
}

func (n *BlockRef) Pos() token.Pos {
	if n == nil || len(n.Args) <= 0 || n.Args[0] == nil {
		return token.NoPos
	}
	return n.Args[0].Pos()
}

func (n *BlockRef) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Args, yield)
}
