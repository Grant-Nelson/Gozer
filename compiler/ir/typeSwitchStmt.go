package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// TypeSwitchStmt is a node that represents a type switch statement.
type TypeSwitchStmt struct {
	Ast    *ast.TypeSwitchStmt // TODO: REMOVE
	Init   Stmt                // initialization statement; or nil
	Assign Stmt                // x := y.(type) or y.(type)
	Body   []*CaseClause
}

var _ Stmt = (*TypeSwitchStmt)(nil)

func (s *TypeSwitchStmt) String() string {
	return fmt.Sprintf("switch %v {\n%s\n}", s.Assign, linesString(s.Body))
}

func (s *TypeSwitchStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*TypeSwitchStmt) StmtNode() {}

func FromTypeSwitchStmt(s *ast.TypeSwitchStmt, c *Converter) *TypeSwitchStmt {
	if s == nil {
		return nil
	}
	return &TypeSwitchStmt{
		Ast:    s,
		Init:   FromStmt(s.Init, c),
		Assign: FromStmt(s.Assign, c),
		Body:   FromCaseClauseSlice(s.Body, c),
	}
}
