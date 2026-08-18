package ir

import (
	"go/token"
	"go/types"
)

type ConstRef struct {
	RefPos token.Pos

	ConstObj *types.Const
}

var (
	_ Expr = (*ConstRef)(nil)
	_ Ref  = (*ConstRef)(nil)
)

func (*ConstRef) ExprNode() {}
func (*ConstRef) RefNode()  {}

func (n *ConstRef) Pos() token.Pos       { return n.RefPos }
func (n *ConstRef) Type() types.Type     { return n.ConstObj.Type() }
func (n *ConstRef) Object() types.Object { return n.ConstObj }
func (n *ConstRef) String() string       { return `ref ` + n.ConstObj.String() }
