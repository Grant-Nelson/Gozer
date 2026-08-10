package ir

import (
	"fmt"
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

func (n *ForStmt) String() string {
	init := ``
	if n.Init != nil {
		init = n.Init.String()
	}
	cond := `true`
	if n.Cond != nil {
		cond = n.Cond.String()
	}
	post := ``
	if n.Post != nil {
		post = n.Post.String()
	}
	var body string
	if len(n.Body) == 1 {
		body = n.Body[0].String()
	} else {
		body = "{\n" + linesString(n.Body) + "\n}"
	}
	return fmt.Sprintf("for (%s; %s; %s) %s", init, cond, post, body)
}

func (n *ForStmt) Pos() token.Pos { return n.ForPos }

func (*ForStmt) StmtNode() {}

func (n *ForStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Cond) && yield(n.Post) && YieldSlice(n.Body, yield)
}
