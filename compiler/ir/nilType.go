package ir

import (
	"go/token"
	"go/types"
)

// NilType is a defined untyped nil value.
type NilType struct {
	Obj types.Object
}

var (
	_ Expr = (*NilType)(nil)
	_ Stmt = (*NilType)(nil)
	_ Decl = (*NilType)(nil)
)

func (n *NilType) ExprNode() {}
func (n *NilType) StmtNode() {}
func (n *NilType) DeclNode() {}

func (n *NilType) Pos() token.Pos       { return n.Obj.Pos() }
func (n *NilType) Type() types.Type     { return n.Obj.Type() }
func (n *NilType) Object() types.Object { return n.Obj }
func (n *NilType) String() string       { return `nil` }
