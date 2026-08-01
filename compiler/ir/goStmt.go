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

var _ Stmt = (*GoStmt)(nil)

func (s *GoStmt) String() string { return fmt.Sprintf(`go %s`, nodeString(s.Call)) }

func (s *GoStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*GoStmt) StmtNode() {}

func FromGoStmt(s *ast.GoStmt) *GoStmt {
	if s == nil {
		return nil
	}
	return &GoStmt{
		Ast:  s,
		Call: s.Call,
	}
}
