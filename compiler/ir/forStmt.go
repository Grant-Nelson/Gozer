package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ForStmt is a node that represents a for statement.
type ForStmt struct {
	Ast  *ast.ForStmt // TODO: REMOVE
	Init Stmt         // initialization statement; or nil
	Cond ast.Expr     // TODO: REPLACE // condition; or nil
	Post Stmt         // post iteration statement; or nil
	Body []Stmt
}

var _ Stmt = (*ForStmt)(nil)

func (s *ForStmt) String() string {
	return fmt.Sprintf("for %v; %s; %v {\n%v\n}",
		s.Init, nodeString(s.Cond), s.Post, linesString(s.Body))
}

func (s *ForStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*ForStmt) StmtNode() {}

func FromForStmt(s *ast.ForStmt, c *Converter) *ForStmt {
	if s == nil {
		return nil
	}
	return &ForStmt{
		Ast:  s,
		Init: c.FromStmt(s.Init),
		Cond: s.Cond,
		Post: c.FromStmt(s.Post),
		Body: c.ExpandStmt(FromBlockStmt(s.Body, c)),
	}
}
