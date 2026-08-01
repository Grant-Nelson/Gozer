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

var (
	_ Stmt   = (*AssignStmt)(nil)
	_ Parent = (*AssignStmt)(nil)
)

func (n *AssignStmt) String() string {
	def := ` = `
	if n.Define {
		def = ` := `
	}
	return fmt.Sprintf(`%s%s%s`, csvString(n.Lhs), def, csvString(n.Rhs))
}

func (n *AssignStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*AssignStmt) StmtNode() {}

func (n *AssignStmt) Children(yield func(Node) bool) bool {
	return yield(n.Lhs) && yield(n.Rhs)
}

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
