package ir

import (
	"go/ast"
	"go/token"
)

// ExprStmt is a node that represents a (stand-alone) expression
// in a statement list.
type ExprStmt struct {
	Ast ast.Stmt // TODO: REMOVE
	X   Expr
}

var (
	_ Stmt   = (*ExprStmt)(nil)
	_ Parent = (*ExprStmt)(nil)
)

func (n *ExprStmt) String() string { return n.X.String() }

func (n *ExprStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*ExprStmt) StmtNode() {}

func (n *ExprStmt) Children(yield func(Node) bool) {
	_ = yield(n.X)
}
