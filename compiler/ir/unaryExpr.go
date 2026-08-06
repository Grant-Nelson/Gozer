package ir

import (
	"fmt"
	"go/token"
)

// UnaryExpr is a node that represents a unary expression.
type UnaryExpr struct {

	// OpPos is the position of the operator.
	OpPos token.Pos

	// Op is the unary operator token (e.g. +, -, !, ^, &, <-).
	Op token.Token

	// X is the operand.
	X Expr
}

var (
	_ Expr   = (*UnaryExpr)(nil)
	_ Parent = (*UnaryExpr)(nil)
)

func (n *UnaryExpr) Pos() token.Pos { return n.OpPos }

func (*UnaryExpr) ExprNode() {}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf(`%s%s`, n.Op, n.X)
}

func (n *UnaryExpr) Children(yield func(Node) bool) {
	_ = yield(n.X)
}
