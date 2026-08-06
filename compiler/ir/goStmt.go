package ir

import (
	"go/ast"
	"go/token"
)

// A GoStmt node represents a go statement.
type GoStmt struct {
	Ast  *ast.GoStmt // TODO: REMOVE
	Call *CallExpr   // TODO: REPLACE
}

var (
	_ Stmt   = (*GoStmt)(nil)
	_ Parent = (*GoStmt)(nil)
)

func (n *GoStmt) String() string { return `go ` + n.Call.String() }

func (n *GoStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*GoStmt) StmtNode() {}

func (n *GoStmt) Children(yield func(Node) bool) {
	_ = yield(n.Call)
}
