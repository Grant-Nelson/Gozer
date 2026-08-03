package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// An IfStmt node represents an if statement.
type IfStmt struct {
	Ast  *ast.IfStmt // TODO: REMOVE
	Init Stmt        // initialization statement; or nil
	Cond Expr        // condition
	Body []Stmt
	Else []Stmt
}

var (
	_ Stmt   = (*IfStmt)(nil)
	_ Parent = (*IfStmt)(nil)
)

func (n *IfStmt) String() string {
	str := fmt.Sprintf("if %s {\n%s\n}", n.Cond, linesString(n.Body))
	if len(n.Else) > 0 {
		str += fmt.Sprintf(" else {\n%s\n}", linesString(n.Else))
	}
	return str
}

func (n *IfStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*IfStmt) StmtNode() {}

func (n *IfStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Cond) &&
		YieldSlice(n.Body, yield) && YieldSlice(n.Else, yield)
}
