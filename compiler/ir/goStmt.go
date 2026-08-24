package ir

import "go/token"

// A GoStmt node represents a go statement.
type GoStmt struct {
	GoPos token.Pos // position of "go" keyword
	Call  *CallExpr
}

var (
	_ Stmt   = (*GoStmt)(nil)
	_ Parent = (*GoStmt)(nil)
)

func (*GoStmt) StmtNode() {}

func (n *GoStmt) Pos() token.Pos { return n.GoPos }
func (n *GoStmt) String() string { return `go ` + toString(n.Call) }

func (n *GoStmt) Children(yield func(Node) bool) { _ = yield(n.Call) }
