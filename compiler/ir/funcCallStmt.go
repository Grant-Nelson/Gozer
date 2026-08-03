package ir

import (
	"go/ast"
	"go/token"
)

// FuncCallStmt is a statement for jumping to another block or named function.
// This is unique to IR and not a mirror of an AST node.
type FuncCallStmt struct {
	Ast    *ast.CallExpr // TODO: REMOVE
	Fun    Expr          // TODO: REPLACE // function expression
	Args   []Expr        // TODO: REPLACE
	Follow *BlockRef
}

var (
	_ Stmt   = (*FuncCallStmt)(nil)
	_ Parent = (*FuncCallStmt)(nil)
)

func (n *FuncCallStmt) String() string {
	return `call(` + n.Fun.String() + `, ` + csvString(n.Args) + `|` + n.Follow.String() + `)`
}

func (n *FuncCallStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*FuncCallStmt) StmtNode() {}

func (n *FuncCallStmt) Children(yield func(Node) bool) {
	_ = yield(n.Fun) && YieldSlice(n.Args, yield) && yield(n.Follow)
}
