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
)

func (*NilType) ExprNode() {}
func (*NilType) StmtNode() {}

func (n *NilType) Pos() token.Pos       { return n.Obj.Pos() }
func (n *NilType) Type() types.Type     { return n.Obj.Type() }
func (n *NilType) Object() types.Object { return n.Obj }
func (n *NilType) String() string       { return `nil` }
