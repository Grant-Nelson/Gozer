package ir

import (
	"go/token"
	"go/types"
)

type KeyValueExpr struct {
	Key   Expr
	Value Expr
}

var _ Expr = (*KeyValueExpr)(nil)

func (n *KeyValueExpr) ExprNode()        {}
func (n *KeyValueExpr) Pos() token.Pos   { return n.Key.Pos() }
func (n *KeyValueExpr) Type() types.Type { return nil }

func (n *KeyValueExpr) String() string {
	return toString(n.Key) + `: ` + toString(n.Value)
}
