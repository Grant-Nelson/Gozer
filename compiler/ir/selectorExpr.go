package ir

import (
	"go/token"
	"go/types"
)

// SelectorExpr is a node that represents an expression followed by a selector.
type SelectorExpr struct {

	// X is the expression being selected from.
	X Expr

	// Sel is the field selector name.
	Sel string

	// SelPos is the position of the selector name.
	SelPos token.Pos

	// SelType is the type of the selected field.
	SelType types.Type
}

var _ Expr = (*SelectorExpr)(nil)

func (n *SelectorExpr) Pos() token.Pos { return n.SelPos }

func (n *SelectorExpr) Type() types.Type { return n.SelType }

func (n *SelectorExpr) ExprNode() {}

func (n *SelectorExpr) String() string { return n.X.String() + `.` + n.Sel }
