package ir

import (
	"go/token"
	"go/types"
)

// IndexListExpr is a node that represents an expression followed by a list of indices.
type IndexListExpr struct {

	// X is the expression that is being indexed.
	X Expr

	// LeftPos is the position of the "["
	LeftPos token.Pos

	// Indices is the expression for the indices.
	Indices []Expr

	// ResultType is the resulting type after this assert.
	ResultType types.Type
}

var _ Expr = (*IndexListExpr)(nil)

func (n *IndexListExpr) ExprNode()        {}
func (n *IndexListExpr) Pos() token.Pos   { return n.LeftPos }
func (n *IndexListExpr) Type() types.Type { return n.ResultType }

func (n *IndexListExpr) String() string {
	return toString(n.X) + `[` + csvString(n.Indices) + `]`
}
