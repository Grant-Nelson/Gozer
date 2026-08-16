package ir

import (
	"go/token"
	"go/types"
)

type TypeRef struct {
	RefPos token.Pos

	TypeDecl *TypeDecl

	TypeArgs []types.Type

	// Instance is the instance type for the signature of the function.
	// This includes the type arguments for the reference. For example,
	// if the declaration is `Foo[T any]`, this reference may be `Foo[int]`.
	Instance types.Type
}

var (
	_ Expr = (*TypeRef)(nil)
	_ Ref  = (*TypeRef)(nil)
)

func (*TypeRef) ExprNode() {}
func (*TypeRef) RefNode()  {}

func (n *TypeRef) Pos() token.Pos   { return n.RefPos }
func (n *TypeRef) Type() types.Type { return n.Instance }
func (n *TypeRef) Decl() Decl       { return n.TypeDecl }
func (n *TypeRef) String() string   { return n.Instance.String() }
