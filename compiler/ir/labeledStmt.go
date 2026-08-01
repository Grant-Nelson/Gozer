package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// A LabeledStmt node represents a labeled statement.
type LabeledStmt struct {
	Ast   *ast.LabeledStmt // TODO: REMOVE
	Label *ast.Ident       // TODO: REPLACE // same as AST so works with types.Info
	Stmt  Stmt             // TODO: INVERT RELATIONSHIP
}

var _ Stmt = (*LabeledStmt)(nil)

func (s *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", s.Label, s.Stmt) }

func (s *LabeledStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*LabeledStmt) StmtNode() {}

func FromLabeledStmt(s *ast.LabeledStmt, c *Converter) *LabeledStmt {
	if s == nil {
		return nil
	}
	return &LabeledStmt{
		Ast:   s,
		Label: s.Label,
		Stmt:  c.FromStmt(s.Stmt),
	}
}
