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
	_ Expr = (*ImportRef)(nil)
	_ Ref  = (*ImportRef)(nil)
)

func (*ImportRef) ExprNode() {}
func (*ImportRef) RefNode()  {}

func (n *ImportRef) Pos() token.Pos       { return n.RefPos }
func (n *ImportRef) Type() types.Type     { return n.ImportDecl.Type() }
func (n *ImportRef) Object() types.Object { return n.ImportDecl.Object() }
func (n *ImportRef) Decl() *ImportDecl    { return n.ImportDecl }
func (n *ImportRef) String() string       { return toString(n.ImportDecl) }
