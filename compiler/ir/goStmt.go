package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// A GoStmt node represents a go statement.
type GoStmt struct {
	Ast  *ast.GoStmt   // TODO: REMOVE
	Call *ast.CallExpr // TODO: REPLACE
}

var (
	_ Stmt   = (*GoStmt)(nil)
	_ Parent = (*GoStmt)(nil)
)

func (n *GoStmt) String() string { return fmt.Sprintf(`go %s`, nodeString(n.Call)) }

func (n *GoStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*GoStmt) StmtNode() {}

func (n *GoStmt) Children(yield func(Node) bool) bool {
	return yield(n.Call)
}

func FromGoStmt(s *ast.GoStmt) *GoStmt {
	if s == nil {
		return nil
	}
	return &GoStmt{
		Ast:  s,
		Call: s.Call,
	}
}
