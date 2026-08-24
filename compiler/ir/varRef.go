package ir

import (
	"go/token"
	"go/types"
)

type VarRef struct {
	RefPos token.Pos
	VarObj *types.Var
}

var (
	_ Expr = (*VarRef)(nil)
	_ Ref  = (*VarRef)(nil)
)

func (*VarRef) ExprNode() {}
func (*VarRef) RefNode()  {}

func (n *VarRef) Pos() token.Pos       { return n.RefPos }
func (n *VarRef) Type() types.Type     { return n.VarObj.Type() }
func (n *VarRef) Object() types.Object { return n.VarObj }
func (n *VarRef) String() string       { return `ref ` + toString(n.VarObj) }
