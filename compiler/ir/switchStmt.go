package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SwitchStmt is a node that represents an expression switch statement.
type SwitchStmt struct {
	Ast  *ast.SwitchStmt
	Init Stmt     // initialization statement; or nil
	Tag  ast.Expr // tag expression; or nil
	Body []*CaseClause
}

var _ Stmt = (*SwitchStmt)(nil)

func (s *SwitchStmt) String() string {
	return fmt.Sprintf("switch %s {\n%s\n}", nodeString(s.Tag), linesString(s.Body))
}

func (s *SwitchStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*SwitchStmt) StmtNode() {}

func FromSwitchStmt(s *ast.SwitchStmt, c *Converter) *SwitchStmt {
	if s == nil {
		return nil
	}
	return &SwitchStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Tag:  s.Tag,
		Body: FromCaseClauseSlice(s.Body, c),
	}
}
