package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// DeferStmt is a node that represents a defer statement.
type DeferStmt struct {
	Ast  *ast.DeferStmt
	Call *ast.CallExpr
}

var _ Stmt = (*DeferStmt)(nil)

func (s *DeferStmt) String() string { return fmt.Sprintf(`defer %s`, nodeString(s.Call)) }

func (s *DeferStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*DeferStmt) StmtNode() {}

func FromDeferStmt(s *ast.DeferStmt) *DeferStmt {
	if s == nil {
		return nil
	}
	return &DeferStmt{
		Ast:  s,
		Call: s.Call,
	}
}
