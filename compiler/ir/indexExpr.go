package ir

import "go/token"

// IndexExpr is a node that represents an expression followed by an index.
type IndexExpr struct {

	// X is the expression that is being indexed.
	X Expr

	// LeftPos is the position of the "["
	LeftPos token.Pos

	// Index is the expression for the index.
	Index Expr
}

var _ Expr = (*IndexExpr)(nil)

func (n *IndexExpr) Pos() token.Pos { return n.LeftPos }

func (n *IndexExpr) ExprNode() {}

func (n *IndexExpr) String() string { return n.X.String() + `[` + n.Index.String() + `]` }
