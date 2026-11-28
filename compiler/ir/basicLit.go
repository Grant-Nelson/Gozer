package ir

import (
	"go/constant"
	"go/token"
	"go/types"
)

// BasicLit is a node that represents a literal of basic type.
type BasicLit struct {

	// ValuePos is the position of the literal.
	ValuePos token.Pos

	// ValueType is the type of the literal.
	ValueType types.Type

	// Value is the value of the literal.
	Value constant.Value
}

var (
	_ Expr = (*BasicLit)(nil)
	_ Stmt = (*BasicLit)(nil)
)

func (*BasicLit) ExprNode() {}
func (*BasicLit) StmtNode() {}

func (n *BasicLit) Pos() token.Pos   { return n.ValuePos }
func (n *BasicLit) Type() types.Type { return n.ValueType }
func (n *BasicLit) String() string   { return toString(n.Value) }
