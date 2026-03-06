package irc

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/project/enums/typeMod"
)

type (
	// Decl is a node interface for all declarations.
	Decl interface {
		Node
		declNode()
	}

	// ValueDecl is a node represents a declaration in a statement list.
	ValueDecl struct {
		Doc    *ast.CommentGroup // associated documentation; or nil
		TokPos token.Pos         // position of the token for the mod
		Mod    typeMod.TypeMod   // modifiers for this value such as const or var
		Lparen token.Pos         // position of '(', if any
		Specs  []*ValueSpec
		Rparen token.Pos // position of ')', if any
	}

	// ValueSpec is a node represents a constant or variable declaration
	// (ConstSpec or VarSpec production).
	ValueSpec struct {
		Doc     *ast.CommentGroup // associated documentation; or nil
		Names   []Ident           // value names (len(Names) > 0)
		Type    Expr              // value type; or nil
		Values  []Expr            // initial values; maybe nil if var
		Comment *ast.CommentGroup // line comments; or nil
	}
)

var (
	_ Decl = (*ValueDecl)(nil)
	_ Node = (*ValueSpec)(nil)
)

//====[String]==================================================================

func (d *ValueDecl) String() string {
	if len(d.Specs) > 1 {
		return fmt.Sprintf("%s (\n%s\n)", d.Mod.String(), linesString(d.Specs, `  `))
	}
	return fmt.Sprintf("%s (\n%s\n)", d.Mod.String(), d.Specs[0])
}

func (s *ValueSpec) String() string {
	return fmt.Sprintf("%s = %s", csvString(s.Names), csvString(s.Values))
}

//====[Pos]=====================================================================

func (d *ValueDecl) Pos() token.Pos { return d.TokPos }
func (s *ValueSpec) Pos() token.Pos { return s.Names[0].Pos() }

//====[End]=====================================================================

func (d *ValueDecl) End() token.Pos {
	if d.Rparen.IsValid() {
		return d.Rparen + 1
	}
	return d.Specs[len(d.Specs)-1].End()
}

func (s *ValueSpec) End() token.Pos {
	if p := endOfSlice(s.Values); p.IsValid() {
		return p
	}
	if s.Type != nil {
		return s.Type.End()
	}
	return s.Names[len(s.Names)-1].End()
}

//====[declNode]================================================================

func (d *ValueDecl) declNode() {}
