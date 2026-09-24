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

	// Slice3 is true if 3-index slice (2 colons present).
	// If true, Max may be nil if unspecified, e.g. `s[x:y:]`.
	// If false, Max should be nil.
	Slice3 bool

	// ResultType is the resulting type after this assert.
	ResultType types.Type
}

var (
	_ Expr   = (*SliceExpr)(nil)
	_ Parent = (*SliceExpr)(nil)
)

func (n *SliceExpr) ExprNode() {}

func (n *SliceExpr) Pos() token.Pos   { return n.LeftPos }
func (n *SliceExpr) Type() types.Type { return n.ResultType }

func (n *SliceExpr) String() string {
	str := toString(n.X) + `[`
	if n.Low != nil {
		str += toString(n.Low)
	}
	str += `:`
	if n.High != nil {
		str += toString(n.High)
	}
	if n.Slice3 {
		str += `:`
		if n.Max != nil {
			str += toString(n.Max)
		}
	}
	return str + `]`
}

func (n *SliceExpr) ChildCount() int { return 4 }

func (n *SliceExpr) Children(yield func(Node) bool) {
	_ = YieldNode(n.X, yield) &&
		YieldNode(n.Low, yield) &&
		YieldNode(n.High, yield) &&
		YieldNode(n.Max, yield)
}
