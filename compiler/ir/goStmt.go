package ir

import (
	"go/token"
)

// A GoStmt node represents a go statement.
type GoStmt struct {
	GoPos token.Pos // position of "go" keyword
	Call  *CallExpr
}

var (
	_ Stmt   = (*GoStmt)(nil)
	_ Parent = (*GoStmt)(nil)
)

func (n *GoStmt) String() string { return `go ` + n.Call.String() }

func (n *GoStmt) Pos() token.Pos { return n.GoPos }

func (*GoStmt) StmtNode() {}

func (n *GoStmt) Children(yield func(Node) bool) {
	_ = yield(n.Call)
}
