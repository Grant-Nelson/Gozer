package ir

import (
	"go/token"
	"go/types"
)

type ValueRef struct {
	RefPos token.Pos

	ValueDecl *ValueDecl
}

var (
	_ Expr = (*ValueRef)(nil)
	_ Ref  = (*ValueRef)(nil)
)

func (*ValueRef) ExprNode() {}
func (*ValueRef) RefNode()  {}

func (n *ValueRef) Pos() token.Pos   { return n.RefPos }
func (n *ValueRef) Type() types.Type { return n.ValueDecl.TypeAndValue.Type }
func (n *ValueRef) Decl() Decl       { return n.ValueDecl }
func (n *ValueRef) String() string   { return `ref ` + n.ValueDecl.String() }
