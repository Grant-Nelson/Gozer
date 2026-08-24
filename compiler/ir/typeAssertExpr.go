package ir

import (
	"go/token"
	"go/types"
)

// TypeAssertExpr is a node that represents a type assertion expression.
type TypeAssertExpr struct {

	// X is the expression whose type is being asserted.
	X Expr

	// LparenPos is the position of the "(" following the dot.
	LparenPos token.Pos

	// AssertType is the asserted type; nil means type switch `x.(type)`.
	AssertType Expr

	// ResultType is the resulting type after this assert.
	ResultType types.Type
}

var (
	_ Expr   = (*TypeAssertExpr)(nil)
	_ Parent = (*TypeAssertExpr)(nil)
)

func (*TypeAssertExpr) ExprNode() {}

func (n *TypeAssertExpr) Pos() token.Pos   { return n.LparenPos }
func (n *TypeAssertExpr) Type() types.Type { return n.ResultType }

func (n *TypeAssertExpr) String() string {
	typeStr := `type`
	if n.AssertType != nil {
		typeStr = toString(n.AssertType)
	}
	return toString(n.X) + `.(` + typeStr + `)`
}

func (n *TypeAssertExpr) Children(yield func(Node) bool) {
	_ = yield(n.X) && yield(n.AssertType)
}
