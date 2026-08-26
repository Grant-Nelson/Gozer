package ir

import (
	"go/token"
	"go/types"
)

// IndexExpr is a node that represents an expression followed by an index.
type IndexExpr struct {

	// X is the expression that is being indexed.
	X Expr

	// LeftPos is the position of the "["
	LeftPos token.Pos

	// Index is the expression for the index.
	Index Expr

	// ResultType is the type resulting from this index.
	ResultType types.Type
}

var _ Expr = (*IndexExpr)(nil)

func (n *IndexExpr) ExprNode() {}

func (n *IndexExpr) Pos() token.Pos   { return n.LeftPos }
func (n *IndexExpr) Type() types.Type { return n.ResultType }

func (n *IndexExpr) String() string {
	return paren(n.X) + `[` + toString(n.Index) + `]`
}
