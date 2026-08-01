package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// StmtListStmt is a node that represents a braced statement list.
type StmtListStmt struct {
	Ast  *ast.BlockStmt // TODO: REMOVE
	List []Stmt
}

var _ Stmt = (*StmtListStmt)(nil)

func (s *StmtListStmt) String() string { return fmt.Sprintf("{\n%s\n}", linesString(s.List)) }

func (s *StmtListStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*StmtListStmt) StmtNode() {}

func FromBlockStmt(s *ast.BlockStmt, c *Converter) *StmtListStmt {
	if s == nil {
		return nil
	}
	return &StmtListStmt{
		Ast:  s,
		List: FromStmtSlice(s.List, c),
	}
}
