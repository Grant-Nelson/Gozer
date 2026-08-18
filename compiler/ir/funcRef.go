package ir

import (
	"go/token"
	"go/types"
)

// FuncRef is an identifier referencing a function or method.
type FuncRef struct {
	// RefPos is the identifier position
	RefPos token.Pos

	FuncObj *types.Func

	TypeArgs []types.Type

	// Instance is the instance type for the signature of the function.
	// This includes the type arguments for the reference. For example,
	// if the declaration is `Foo[T any]`, this reference may be `Foo[int]`.
	Instance types.Type
}

var (
	_ Expr = (*FuncRef)(nil)
	_ Ref  = (*FuncRef)(nil)
)

func (*FuncRef) ExprNode() {}
func (*FuncRef) RefNode()  {}

func (n *FuncRef) Pos() token.Pos       { return n.RefPos }
func (n *FuncRef) Type() types.Type     { return n.Instance }
func (n *FuncRef) Object() types.Object { return n.FuncObj }
func (n *FuncRef) String() string       { return n.Instance.String() }
