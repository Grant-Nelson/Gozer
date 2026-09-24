package ir

import (
	"go/token"
	"go/types"
)

// TODO: See if there is a better way to store type initializations

type KeyValueExpr struct {
	Key   Expr
	Value Expr
}

var (
	_ Expr   = (*KeyValueExpr)(nil)
	_ Parent = (*KeyValueExpr)(nil)
)

func (n *KeyValueExpr) ExprNode()        {}
func (n *KeyValueExpr) Pos() token.Pos   { return n.Key.Pos() }
func (n *KeyValueExpr) Type() types.Type { return nil }

func (n *KeyValueExpr) String() string {
	return toString(n.Key) + `: ` + toString(n.Value)
}

func (n *KeyValueExpr) ChildCount() int { return 2 }

func (n *KeyValueExpr) Children(yield func(Node) bool) {
	_ = YieldNode(n.Key, yield) &&
		YieldNode(n.Value, yield)
}
