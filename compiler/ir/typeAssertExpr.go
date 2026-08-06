package ir

import (
	"fmt"
	"go/token"
)

// TypeAssertExpr is a node that represents a type assertion expression.
type TypeAssertExpr struct {

	// X is the expression whose type is being asserted.
	X Expr

	// LparenPos is the position of the "(" following the dot.
	LparenPos token.Pos

	// Type is the asserted type; nil means type switch `x.(type)`.
	Type Expr
}

var (
	_ Expr   = (*TypeAssertExpr)(nil)
	_ Parent = (*TypeAssertExpr)(nil)
)

func (n *TypeAssertExpr) Pos() token.Pos { return n.LparenPos }

func (*TypeAssertExpr) ExprNode() {}

func (n *TypeAssertExpr) String() string {
	typeStr := `type`
	if n.Type != nil {
		typeStr = n.Type.String()
	}
	return fmt.Sprintf(`%s.(%s)`, n.X, typeStr)
}

func (n *TypeAssertExpr) Children(yield func(Node) bool) {
	_ = yield(n.X) && yield(n.Type)
}
