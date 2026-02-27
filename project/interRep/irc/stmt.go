package irc

import "go/token"

type (
	Stmt interface {
		// Pos is the position for this statement in the source file.
		// This may be [token.NoPos] if not specifically in the source.
		// This may be able to be used to look up the [ast.Node] that
		// became this IRC statement.
		Pos() token.Pos

		// stmt is an empty method used to compile time type check that
		// only statements are used for this interface.
		stmt()
	}

	// ExprStmt is a statement for an expression.
	ExprStmt struct {
		Expr Expr
	}

	ReturnStmt struct {
		Return  token.Pos // the position of "return" keyword
		Results Expr      // the result expression; or nil
	}
)

func (s *ExprStmt) Pos() token.Pos   { return s.Expr.Pos() }
func (s *ReturnStmt) Pos() token.Pos { return s.Return }

func (*ExprStmt) stmt()   {}
func (*ReturnStmt) stmt() {}
