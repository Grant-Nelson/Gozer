package ir

import (
	"fmt"
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/unaryOp"
)

// UnaryExpr is a node that represents a unary expression.
type UnaryExpr struct {

	// OpPos is the position of the operator.
	OpPos token.Pos

	// Op is the unary operator token (e.g. +, -, !, ^, &, <-).
	Op unaryOp.UnaryOp

	// X is the operand.
	X Expr

	// ResultType is the type resulting from this unary expression.
	ResultType types.Type
}

var (
	_ Expr   = (*UnaryExpr)(nil)
	_ Parent = (*UnaryExpr)(nil)
)

func (n *UnaryExpr) Pos() token.Pos { return n.OpPos }

func (n *UnaryExpr) Type() types.Type { return n.ResultType }

func (*UnaryExpr) ExprNode() {}

func (n *UnaryExpr) String() string {
	return fmt.Sprintf(`(%s%s)`, n.Op, n.X)
}

func (n *UnaryExpr) Children(yield func(Node) bool) {
	_ = yield(n.X)
}
