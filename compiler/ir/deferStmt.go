package ir

import (
	"go/token"
)

// DeferStmt is a node that represents a defer statement.
type DeferStmt struct {
	// Defer is the position of "defer" keyword
	Defer token.Pos

	// Call is the function call expression that is being deferred.
	Call *CallExpr
}

var (
	_ Stmt   = (*DeferStmt)(nil)
	_ Parent = (*DeferStmt)(nil)
)

func (n *DeferStmt) String() string { return `defer ` + n.Call.String() }

func (n *DeferStmt) Pos() token.Pos { return n.Defer }

func (*DeferStmt) StmtNode() {}

func (n *DeferStmt) Children(yield func(Node) bool) {
	_ = yield(n.Call)
}
