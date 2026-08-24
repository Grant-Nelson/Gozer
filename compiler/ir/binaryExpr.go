package ir

import (
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/binaryOp"
)

// BinaryExpr is a node that represents a binary expression.
type BinaryExpr struct {

	// X is the left operand.
	X Expr

	// OpPos is the position of the operator.
	OpPos token.Pos

	// Op is the binary operator token (e.g. +, -, *, /, ==, <, &&, ||).
	Op binaryOp.BinaryOp

	// Y is the right operand.
	Y Expr

	// ResultType is the type resulting from this binary expression.
	ResultType types.Type
}

var (
	_ Expr   = (*BinaryExpr)(nil)
	_ Parent = (*BinaryExpr)(nil)
)

func (*BinaryExpr) ExprNode() {}

func (n *BinaryExpr) Pos() token.Pos   { return n.OpPos }
func (n *BinaryExpr) Type() types.Type { return n.ResultType }

func (n *BinaryExpr) String() string {
	lhs := `_`
	if n.X != nil {
		lhs = `(` + toString(n.X) + `)`
	}
	rhs := `_`
	if n.Y != nil {
		rhs = `(` + toString(n.Y) + `)`
	}
	return lhs + toString(n.Op) + rhs
}

func (n *BinaryExpr) Children(yield func(Node) bool) {
	_ = yield(n.X) && yield(n.Y)
}
