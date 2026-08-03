package ir

import (
	"go/token"
)

// ExprStmt is a node that represents a (stand-alone) expression
// in a statement list.
type ExprStmt struct {

	// X is the expression in the statement list.
	X Expr
}

var (
	_ Stmt   = (*ExprStmt)(nil)
	_ Parent = (*ExprStmt)(nil)
)

func (n *ExprStmt) String() string { return n.X.String() }

func (n *ExprStmt) Pos() token.Pos { return n.X.Pos() }

func (*ExprStmt) StmtNode() {}

func (n *ExprStmt) Children(yield func(Node) bool) {
	_ = yield(n.X)
}
