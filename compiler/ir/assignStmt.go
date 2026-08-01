package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// AssignStmt is a node that represents an assignment or
// a short variable declaration.
type AssignStmt struct {
	Ast    *ast.AssignStmt // TODO: REMOVE
	Lhs    []ast.Expr      // TODO: REPLACE
	Define bool
	Rhs    []ast.Expr // TODO: REPLACE
}

var _ Stmt = (*AssignStmt)(nil)

func (s *AssignStmt) String() string {
	def := ` = `
	if s.Define {
		def = ` := `
	}
	return fmt.Sprintf(`%s%s%s`, csvString(s.Lhs), def, csvString(s.Rhs))
}

func (s *AssignStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*AssignStmt) StmtNode() {}

func FromAssignStmt(s *ast.AssignStmt) *AssignStmt {
	if s == nil {
		return nil
	}
	return &AssignStmt{
		Ast:    s,
		Lhs:    s.Lhs,
		Define: s.Tok == token.DEFINE,
		Rhs:    s.Rhs,
	}
}
