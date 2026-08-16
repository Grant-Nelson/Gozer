package ir

import (
	"go/token"
	"go/types"
)

type ImportRef struct {
	RefPos token.Pos

	ImportDecl *ImportDecl
}

var (
	_ Expr = (*ConstRef)(nil)
	_ Ref  = (*ConstRef)(nil)
)

func (*ImportRef) ExprNode() {}
func (*ImportRef) RefNode()  {}

func (n *ImportRef) Pos() token.Pos   { return n.RefPos }
func (n *ImportRef) Type() types.Type { return n.ImportDecl.Type() }
func (n *ImportRef) Decl() Decl       { return n.ImportDecl }
func (n *ImportRef) String() string   { return n.ImportDecl.String() }
