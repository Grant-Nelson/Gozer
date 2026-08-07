package ir

import (
	"go/token"
	"go/types"
)

// SliceExpr is a node that represents an expression followed by slice indices.
type SliceExpr struct {

	// X is the expression that is being sliced.
	X Expr

	// LeftPos is the position of the "["
	LeftPos token.Pos

	// Low is the begin of slice range; or nil
	Low Expr

	// High is the end of slice range; or nil
	High Expr

	// Max is the maximum capacity of slice; or nil
	Max Expr

	// Slice3 is true if 3-index slice (2 colons present)
	Slice3 bool

	// ResultType is the resulting type after this assert.
	ResultType types.Type
}

var _ Expr = (*SliceExpr)(nil)

func (n *SliceExpr) Pos() token.Pos { return n.LeftPos }

func (n *SliceExpr) Type() types.Type { return n.ResultType }

func (n *SliceExpr) ExprNode() {}

func (n *SliceExpr) String() string {
	str := n.X.String() + `[`
	if n.Low != nil {
		str += n.Low.String()
	}
	str += `:`
	if n.High != nil {
		str += n.High.String()
	}
	if n.Slice3 {
		str += `:`
		if n.Max != nil {
			str += n.Max.String()
		}
	}
	return str + `]`
}
