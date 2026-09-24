package ir

import (
	"go/token"
	"go/types"
)

type TypeLit struct {
	TypePos token.Pos
	TypeRef types.Type
	Values  []Expr
}

var (
	_ Expr   = (*TypeLit)(nil)
	_ Parent = (*TypeLit)(nil)
)

func (n *TypeLit) ExprNode() {}

func (n *TypeLit) Pos() token.Pos   { return n.TypePos }
func (n *TypeLit) Type() types.Type { return n.TypeRef }

func (n *TypeLit) String() string {
	return toString(n.TypeRef) + bodyString(n.Values)
}
func (n *TypeLit) ChildCount() int { return len(n.Values) }

func (n *TypeLit) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Values, yield)
}
