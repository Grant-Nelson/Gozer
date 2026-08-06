package ir

import (
	"go/token"
	"go/types"
)

type TypeExpr struct {
	TypePos      token.Pos
	TypeAndValue types.TypeAndValue
}

var _ Expr = (*TypeExpr)(nil)

func (n *TypeExpr) Pos() token.Pos { return n.TypePos }

func (n *TypeExpr) ExprNode() {}

func (n *TypeExpr) String() string {
	if n.TypeAndValue.Value != nil {
		return n.TypeAndValue.Value.String()
	}
	return n.TypeAndValue.Type.String()
}
