package ir

import (
	"fmt"
	"go/token"
)

// BinaryExpr is a node that represents a binary expression.
type BinaryExpr struct {

	// X is the left operand.
	X Expr

	// OpPos is the position of the operator.
	OpPos token.Pos

	// Op is the binary operator token (e.g. +, -, *, /, ==, <, &&, ||).
	Op token.Token

	// Y is the right operand.
	Y Expr
}

var (
	_ Expr   = (*BinaryExpr)(nil)
	_ Parent = (*BinaryExpr)(nil)
)

func (n *BinaryExpr) Pos() token.Pos { return n.OpPos }

func (*BinaryExpr) ExprNode() {}

func (n *BinaryExpr) String() string {
	return fmt.Sprintf(`%s %s %s`, n.X, n.Op, n.Y)
}

func (n *BinaryExpr) Children(yield func(Node) bool) {
	_ = yield(n.X) && yield(n.Y)
}
