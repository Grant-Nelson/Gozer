package ir

import (
	"go/token"
	"go/types"
)

type BuiltinRef struct {
	RefPos  token.Pos
	Builtin *types.Builtin
}

var (
	_ Expr = (*BuiltinRef)(nil)
	_ Ref  = (*BuiltinRef)(nil)
)

func (*BuiltinRef) ExprNode() {}
func (*BuiltinRef) RefNode()  {}

func (n *BuiltinRef) Pos() token.Pos       { return n.RefPos }
func (n *BuiltinRef) Type() types.Type     { return n.Builtin.Type() }
func (n *BuiltinRef) Object() types.Object { return n.Builtin }
func (n *BuiltinRef) String() string       { return toString(n.Builtin) }
