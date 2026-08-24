package ir

import (
	"go/token"
	"go/types"
)

type TypeExpr struct {
	TypePos      token.Pos
	TypeAndValue *types.TypeAndValue
}

var _ Expr = (*TypeExpr)(nil)

func (n *TypeExpr) ExprNode() {}

func (n *TypeExpr) Pos() token.Pos   { return n.TypePos }
func (n *TypeExpr) Type() types.Type { return n.TypeAndValue.Type }

func (n *TypeExpr) String() string {
	if n.TypeAndValue.Value != nil {
		return toString(n.TypeAndValue.Value)
	}
	return toString(n.TypeAndValue.Type)
}
