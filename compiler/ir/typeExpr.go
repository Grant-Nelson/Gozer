package ir

import (
	"go/token"
	"go/types"
)

type TypeExpr struct {
	TypePos token.Pos
	TypeRef types.Type
}

var (
	_ Expr = (*TypeExpr)(nil)
	_ Stmt = (*TypeExpr)(nil)
)

func (*TypeExpr) ExprNode() {}
func (*TypeExpr) StmtNode() {}

func (n *TypeExpr) Pos() token.Pos   { return n.TypePos }
func (n *TypeExpr) Type() types.Type { return n.TypeRef }

func (n *TypeExpr) String() string { return toString(n.TypeRef) }
