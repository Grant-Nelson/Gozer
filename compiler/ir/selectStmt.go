package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SelectStmt is a node that represents a select statement.
type SelectStmt struct {
	Ast  *ast.SelectStmt // TODO: REMOVE
	Body []*CommClause
}

var _ Stmt = (*SelectStmt)(nil)

func (s *SelectStmt) String() string {
	return fmt.Sprintf("select {\n%s\n}", linesString(s.Body))
}

func (s *SelectStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*SelectStmt) StmtNode() {}

func FromSelectStmt(s *ast.SelectStmt, c *Converter) *SelectStmt {
	if s == nil {
		return nil
	}
	return &SelectStmt{
		Ast:  s,
		Body: FromCommClauseSlice(s.Body, c),
	}
}
