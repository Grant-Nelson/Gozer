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

var (
	_ Stmt   = (*SelectStmt)(nil)
	_ Parent = (*SelectStmt)(nil)
)

func (n *SelectStmt) String() string {
	return fmt.Sprintf("select {\n%s\n}", linesString(n.Body))
}

func (n *SelectStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*SelectStmt) StmtNode() {}

func (n *SelectStmt) Children(yield func(Node) bool) bool {
	return YieldSlice(n.Body, yield)
}

func FromSelectStmt(s *ast.SelectStmt, c *Converter) *SelectStmt {
	if s == nil {
		return nil
	}
	return &SelectStmt{
		Ast:  s,
		Body: FromCommClauseSlice(s.Body, c),
	}
}
