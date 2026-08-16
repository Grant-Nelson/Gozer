package ir

import (
	"go/token"
	"go/types"
)

type ConstRef struct {
	RefPos token.Pos

	ConstDecl *ConstDecl
}

var (
	_ Expr = (*ConstRef)(nil)
	_ Ref  = (*ConstRef)(nil)
)

func (*ConstRef) ExprNode() {}
func (*ConstRef) RefNode()  {}

func (n *ConstRef) Pos() token.Pos   { return n.RefPos }
func (n *ConstRef) Type() types.Type { return n.ConstDecl.Type() }
func (n *ConstRef) Decl() Decl       { return n.ConstDecl }
func (n *ConstRef) String() string   { return `ref ` + n.ConstDecl.String() }
