package ir

import (
	"fmt"
	"go/token"
)

// CallExpr is a node that represents a function or method call expression.
type CallExpr struct {

	// Fun is the function or method being called.
	Fun Expr

	// LparenPos is the position of the "(".
	LparenPos token.Pos

	// Args are the arguments to the call.
	Args []Expr

	// Ellipsis is the position of "..." if the last argument is variadic,
	// or token.NoPos otherwise.
	Ellipsis token.Pos
}

var (
	_ Expr   = (*CallExpr)(nil)
	_ Parent = (*CallExpr)(nil)
)

func (n *CallExpr) Pos() token.Pos { return n.LparenPos }

func (*CallExpr) ExprNode() {}

func (n *CallExpr) String() string {
	ellipsis := ``
	if n.Ellipsis.IsValid() {
		ellipsis = `...`
	}
	return fmt.Sprintf(`%s(%s%s)`, n.Fun, csvString(n.Args), ellipsis)
}

func (n *CallExpr) Children(yield func(Node) bool) {
	_ = yield(n.Fun) && YieldSlice(n.Args, yield)
}
