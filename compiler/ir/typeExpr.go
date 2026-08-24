package ir

import (
	"go/token"
	"go/types"
)

type TypeExpr struct {
	TypePos token.Pos
	TypeRef types.Type
}

var _ Expr = (*TypeExpr)(nil)

func (n *TypeExpr) ExprNode() {}

func (n *TypeExpr) Pos() token.Pos   { return n.TypePos }
func (n *TypeExpr) Type() types.Type { return n.TypeRef }

func (n *TypeExpr) String() string { return toString(n.TypeRef) }
