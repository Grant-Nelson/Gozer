package ir

import (
	"go/token"
	"go/types"
)

// BasicLit is a node that represents a literal of basic type.
type BasicLit struct {

	// ValuePos is the position of the literal.
	ValuePos token.Pos

	// TypeAndValue is the type and value of the literal.
	TypeAndValue types.TypeAndValue
}

var _ Expr = (*BasicLit)(nil)

func (n *BasicLit) Pos() token.Pos { return n.ValuePos }

func (n *BasicLit) Type() types.Type { return n.TypeAndValue.Type }

func (n *BasicLit) ExprNode() {}

func (n *BasicLit) String() string { return n.TypeAndValue.Value.String() }
