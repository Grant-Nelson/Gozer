package ir

import "go/token"

// SelectorExpr is a node that represents an expression followed by a selector.
type SelectorExpr struct {

	// X is the expression being selected from.
	X Expr

	// Sel is the field selector
	Sel *Ident
}

var _ Expr = (*SelectorExpr)(nil)

func (n *SelectorExpr) Pos() token.Pos { return n.Sel.NamePos }

func (n *SelectorExpr) ExprNode() {}

func (n *SelectorExpr) String() string { return n.X.String() + n.Sel.String() }
