package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// StmtListStmt is a node that represents a braced statement list.
type StmtListStmt struct {
	Ast  *ast.BlockStmt // TODO: REMOVE
	List []Stmt
}

var (
	_ Stmt   = (*StmtListStmt)(nil)
	_ Parent = (*StmtListStmt)(nil)
)

func (n *StmtListStmt) String() string { return fmt.Sprintf("{\n%s\n}", linesString(n.List)) }

func (n *StmtListStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*StmtListStmt) StmtNode() {}

func (n *StmtListStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.List, yield)
}
