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

var (
	_ Stmt   = (*DeferStmt)(nil)
	_ Parent = (*DeferStmt)(nil)
)

func (n *DeferStmt) String() string { return fmt.Sprintf(`defer %s`, nodeString(n.Call)) }

func (n *DeferStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*DeferStmt) StmtNode() {}

func (n *DeferStmt) Children(yield func(Node) bool) bool {
	return yield(n.Call)
}

func FromDeferStmt(s *ast.DeferStmt) *DeferStmt {
	if s == nil {
		return nil
	}
	return &DeferStmt{
		Ast:  s,
		Call: s.Call,
	}
}
