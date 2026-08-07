package ir

import (
	"go/token"
	"go/types"
)

// Ident is a node that represents an identifier.
type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name

	TypeAndValue types.TypeAndValue
}

var _ Expr = (*Ident)(nil)

func (n *Ident) Pos() token.Pos { return n.NamePos }

func (n *Ident) Type() types.Type { return n.TypeAndValue.Type }

func (*Ident) ExprNode() {}

func (n *Ident) String() string { return n.Name }
