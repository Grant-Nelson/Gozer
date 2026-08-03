package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ForStmt is a node that represents a for statement.
type ForStmt struct {
	Ast  *ast.ForStmt // TODO: REMOVE
	Init Stmt         // initialization statement; or nil
	Cond Expr         // condition; or nil
	Post Stmt         // post iteration statement; or nil
	Body []Stmt
}

var (
	_ Stmt   = (*ForStmt)(nil)
	_ Parent = (*ForStmt)(nil)
)

func (n *ForStmt) String() string {
	return fmt.Sprintf("for %v; %s; %v {\n%v\n}",
		n.Init, n.Cond, n.Post, linesString(n.Body))
}

func (n *ForStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*ForStmt) StmtNode() {}

func (n *ForStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Cond) && yield(n.Post) && YieldSlice(n.Body, yield)
}
