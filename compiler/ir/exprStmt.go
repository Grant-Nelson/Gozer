package ir

import (
	"go/ast"
	"go/token"
)

// ExprStmt is a node that represents a (stand-alone) expression
// in a statement list.
type ExprStmt struct {
	Ast ast.Stmt // TODO: REMOVE
	X   ast.Expr // TODO: REPLACE
}

var (
	_ Stmt   = (*ExprStmt)(nil)
	_ Parent = (*ExprStmt)(nil)
)

func (n *ExprStmt) String() string { return nodeString(n.X) }

func (n *ExprStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*ExprStmt) StmtNode() {}

func (n *ExprStmt) Children(yield func(Node) bool) bool {
	return yield(n.X)
}

func FromExprStmt(s *ast.ExprStmt) *ExprStmt {
	if s == nil {
		return nil
	}
	return &ExprStmt{
		Ast: s,
		X:   s.X,
	}
}

func FromIncDecStmt(s *ast.IncDecStmt, c *Converter) *ExprStmt {
	if s == nil {
		return nil
	}
	exp := &ast.UnaryExpr{OpPos: s.Pos(), Op: s.Tok, X: s.X}
	c.Info.Types[exp] = c.Info.Types[s.X]
	return &ExprStmt{
		Ast: s,
		X:   exp,
	}
}
