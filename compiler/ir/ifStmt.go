package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// An IfStmt node represents an if statement.
type IfStmt struct {
	Ast  *ast.IfStmt // TODO: REMOVE
	Init Stmt        // initialization statement; or nil
	Cond ast.Expr    // TODO: REPLACE // condition
	Body []Stmt
	Else []Stmt
}

var _ Stmt = (*IfStmt)(nil)

func (s *IfStmt) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}", nodeString(s.Cond), linesString(s.Body))
	if len(s.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", linesString(s.Else))
	}
	return str
}

func (s *IfStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*IfStmt) StmtNode() {}

func FromIfStmt(s *ast.IfStmt, c *Converter) *IfStmt {
	if s == nil {
		return nil
	}
	return &IfStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: s.Cond,
		Body: c.ExpandStmt(FromBlockStmt(s.Body, c)),
		Else: c.ExpandStmt(c.FromStmt(s.Else)),
	}
}
