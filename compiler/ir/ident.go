package ir

import "go/token"

// Ident is a node that represents an identifier.
type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name

	// TODO: Add Additional Info
}

var _ Expr = (*Ident)(nil)

func (n *Ident) Pos() token.Pos { return n.NamePos }

func (*Ident) ExprNode() {}

func (n *Ident) String() string { return n.Name }
