package ir

import "go/token"

// GotoBlockStmt is a flow control statement for a goto jump to another block.
// This is unique to IR and not a mirror of an AST node.
type GotoBlockStmt struct {
	SrcPos token.Pos
	Block  *BlockRef
}

var (
	_ Stmt     = (*GotoBlockStmt)(nil)
	_ FlowCtrl = (*GotoBlockStmt)(nil)
	_ Parent   = (*GotoBlockStmt)(nil)
)

func (*GotoBlockStmt) StmtNode()     {}
func (*GotoBlockStmt) FlowCtrlNode() {}

func (n *GotoBlockStmt) Pos() token.Pos { return n.SrcPos }
func (n *GotoBlockStmt) String() string { return `goto(` + toString(n.Block) + `)` }

func (n *GotoBlockStmt) Children(yield func(Node) bool) { _ = yield(n.Block) }

func NewGotoBlockStmt(pos token.Pos, nextBlk *Block, args ...Expr) *GotoBlockStmt {
	return &GotoBlockStmt{
		SrcPos: pos,
		Block: &BlockRef{
			Block: nextBlk,
			Args:  args,
		},
	}
}
