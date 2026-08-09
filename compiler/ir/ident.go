package ir

import (
	"go/token"
	"go/types"
)

// Ident is a node that represents an identifier.
type Ident struct {
	NamePos token.Pos // identifier position
	Name    string    // identifier name

	TypeAndValue *types.TypeAndValue
	Instance     *types.Instance
	Def          types.Object
	Use          types.Object
}

var _ Expr = (*Ident)(nil)

func (n *Ident) Pos() token.Pos { return n.NamePos }

func (n *Ident) Type() types.Type {
	if n.TypeAndValue != nil && n.TypeAndValue.Type != nil {
		return n.TypeAndValue.Type
	}
	if n.Instance != nil && n.Instance.Type != nil {
		return n.Instance.Type
	}
	if n.Def != nil && n.Def.Type() != nil {
		return n.Def.Type()
	}
	if n.Use != nil && n.Use.Type() != nil {
		return n.Use.Type()
	}
	return nil
}

func (*Ident) ExprNode() {}

func (n *Ident) String() string { return n.Name }

func (n *Ident) DetailedString() string {
	if n.TypeAndValue != nil && n.TypeAndValue.Type != nil {
		return n.Name + ` ` + n.TypeAndValue.Type.String()
	}
	if n.Instance != nil && n.Instance.Type != nil {
		return n.Name + ` (instance)` + n.Instance.Type.String()
	}
	if n.Def != nil && n.Def.Type() != nil {
		return n.Name + ` (def)` + n.Def.Type().String()
	}
	if n.Use != nil && n.Use.Type() != nil {
		return n.Name + ` (use)` + n.Use.Type().String()
	}
	return n.Name + ` (unknown type)`
}
