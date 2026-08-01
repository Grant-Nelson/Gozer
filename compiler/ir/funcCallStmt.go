package ir

import (
	"go/ast"
	"go/token"
)

// FuncCallStmt is a statement for jumping to another block or named function.
// This is unique to IR and not a mirror of an AST node.
type FuncCallStmt struct {
	Ast    *ast.CallExpr // TODO: REMOVE
	Fun    ast.Expr      // TODO: REPLACE // function expression
	Args   []ast.Expr    // TODO: REPLACE
	Follow *BlockRef
}

var _ Stmt = (*FuncCallStmt)(nil)

func (s *FuncCallStmt) String() string {
	return `call(` + nodeString(s.Fun) + `, ` + csvString(s.Args) + `|` + s.Follow.String() + `)`
}

func (s *FuncCallStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*FuncCallStmt) StmtNode() {}
