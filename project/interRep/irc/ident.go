package irc

import (
	"go/token"
	"go/types"
)

type (
	Ident interface {
		Expr

		// identNode is an empty method used to compile time type check that
		// only identifier duck-type to this interface.
		identNode()
	}

	DefIdent struct {
		NamePos token.Pos // identifier position
		Name    string    // identifier name

		// Object is the object defined by this identifier (including
		// package names, dots "." of dot-imports, and blank "_" identifiers).
		// For an embedded field, this is the field *Var it defines.
		Object types.Object
	}

	// TODO: Add more identifiers
)

var (
	_ Ident = (*DefIdent)(nil)
)

// IsExported reports whether id starts with an upper-case letter.
func (i *DefIdent) IsExported() bool { return token.IsExported(i.Name) }

//====[String]==================================================================

func (i *DefIdent) String() string {
	if i != nil {
		return i.Name
	}
	return `<nil>`
}

//====[Pos]=====================================================================

func (i *DefIdent) Pos() token.Pos { return i.NamePos }

//====[End]=====================================================================

func (i *DefIdent) End() token.Pos { return token.Pos(int(i.NamePos) + len(i.Name)) }

//====[ResultType]==============================================================

func (i *DefIdent) ResultType() types.Type { return i.Object.Type() }

//====[declNode]================================================================

func (i *DefIdent) exprNode() {}

//====[declNode]================================================================

func (i *DefIdent) identNode() {}
