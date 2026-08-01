package ir

import (
	"go/ast"
	"go/token"
)

// GotoBlockStmt is a flow control statement for a goto jump to another block.
// This is unique to IR and not a mirror of an AST node.
type GotoBlockStmt struct {
	SrcPos token.Pos
	Block  *BlockRef
}

var (
	_ Stmt     = (*GotoBlockStmt)(nil)
	_ FlowCtrl = (*GotoBlockStmt)(nil)
)

func (s *GotoBlockStmt) String() string {
	return `goto(` + s.Block.String() + `)`
}

func (s *GotoBlockStmt) Pos() token.Pos { return s.SrcPos }

func (*GotoBlockStmt) StmtNode()     {}
func (*GotoBlockStmt) FlowCtrlNode() {}

func NewGotoBlockStmt(pos token.Pos, nextBlk *Block, args ...ast.Expr) *GotoBlockStmt {
	return &GotoBlockStmt{
		SrcPos: pos,
		Block: &BlockRef{
			Block: nextBlk,
			Args:  args,
		},
	}
}
