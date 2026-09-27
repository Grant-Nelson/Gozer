package ir

import (
	"go/token"
	"go/types"
)

type ImportRef struct {
	RefPos token.Pos

	ImportDef *ImportDef
}

var (
	_ Expr = (*ImportRef)(nil)
	_ Stmt = (*ImportRef)(nil)
	_ Ref  = (*ImportRef)(nil)
)

func (*ImportRef) ExprNode() {}
func (*ImportRef) StmtNode() {}
func (*ImportRef) RefNode()  {}

func (n *ImportRef) Pos() token.Pos       { return n.RefPos }
func (n *ImportRef) Type() types.Type     { return n.ImportDef.Type() }
func (n *ImportRef) Object() types.Object { return n.ImportDef.Object() }
func (n *ImportRef) Decl() *ImportDef     { return n.ImportDef }
func (n *ImportRef) String() string       { return toString(n.ImportDef) }
