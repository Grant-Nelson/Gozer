package ir

import (
	"go/token"
	"go/types"
)

// CallExpr is a node that represents a function or method call expression.
type CallExpr struct {
	// Fun is the function or method being called.
	Fun Expr

	// LeftParenPos is the position of the "(".
	LeftParenPos token.Pos

	// Args are the arguments to the call.
	Args []Expr

	// Expanded indicates if there is a "..." for the last variadic argument
	// to expand zero or more values.
	Expanded bool

	// ResultType is the type that is returned from this call.
	ResultType types.Type

	Follow *BlockRef
}

var (
	_ Expr   = (*CallExpr)(nil)
	_ Parent = (*CallExpr)(nil)
)

func (*CallExpr) ExprNode() {}

func (n *CallExpr) Pos() token.Pos   { return n.LeftParenPos }
func (n *CallExpr) Type() types.Type { return n.ResultType }

func (n *CallExpr) String() string {
	ellipsis := ``
	if n.Expanded {
		ellipsis = `...`
	}
	return toString(n.Fun) + `(` + csvString(n.Args) + ellipsis + `)`
}

func (n *CallExpr) Children(yield func(Node) bool) {
	_ = yield(n.Fun) && YieldSlice(n.Args, yield) && yield(n.Follow)
}
