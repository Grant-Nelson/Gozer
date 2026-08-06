package ir

import (
	"go/constant"
	"go/token"
)

// BasicLit is a node that represents a literal of basic type.
type BasicLit struct {

	// ValuePos is the position of the literal.
	ValuePos token.Pos

	// Value is the literal value.
	Value constant.Value
}

var _ Expr = (*BasicLit)(nil)

func (n *BasicLit) Pos() token.Pos { return n.ValuePos }

func (n *BasicLit) ExprNode() {}

func (n *BasicLit) String() string { return n.Value.String() }
