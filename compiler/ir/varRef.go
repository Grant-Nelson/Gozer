package ir

import (
	"go/token"
	"go/types"
)

type VarRef struct {
	RefPos token.Pos

	VarDecl *VarDecl
}

var (
	_ Expr = (*VarRef)(nil)
	_ Ref  = (*VarRef)(nil)
)

func (*VarRef) ExprNode() {}
func (*VarRef) RefNode()  {}

func (n *VarRef) Pos() token.Pos   { return n.RefPos }
func (n *VarRef) Type() types.Type { return n.VarDecl.Type() }
func (n *VarRef) Decl() Decl       { return n.VarDecl }
func (n *VarRef) String() string   { return `ref ` + n.VarDecl.String() }
