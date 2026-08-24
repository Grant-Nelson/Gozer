package ir

import (
	"go/token"
)

// ForStmt is a node that represents a for statement.
type ForStmt struct {
	// ForPos is the position of the `for` keyword.
	ForPos token.Pos

	// Init is the initialization statement; or nil.
	Init Stmt

	// Cond is the condition; or nil.
	Cond Expr

	// Post is the post iteration statement; or nil.
	Post Stmt

	// Body is the statements in the body of the loop.
	Body []Stmt
}

var (
	_ Stmt   = (*ForStmt)(nil)
	_ Parent = (*ForStmt)(nil)
)

func (*ForStmt) StmtNode() {}

func (n *ForStmt) Pos() token.Pos { return n.ForPos }

func (n *ForStmt) String() string {
	if n.IsWhile() {
		return `while (` + emptyZeroOrString(n.Cond) + `)` + bodyString(n.Body)
	}
	return `for (` + emptyZeroOrString(n.Init) + `; ` + emptyZeroOrString(n.Cond) +
		`; ` + emptyZeroOrString(n.Post) + `)` + bodyString(n.Body)
}

func (n *ForStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Cond) && yield(n.Post) && YieldSlice(n.Body, yield)
}

func (n *ForStmt) IsWhile() bool {
	return n.Init == nil && n.Post == nil
}
