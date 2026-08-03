package ir

import (
	"go/ast"
	"go/token"
)

// TODO: REPLACE WITH statements specific for CONST, TYPE, and VAR
//
// DeclStmt is a node that represents a declaration in a statement list.
type DeclStmt struct {
	Ast  *ast.DeclStmt // TODO: REMOVE
	Decl ast.Decl      // TODO: REPLACE // *GenDecl with CONST, TYPE, or VAR token
}

var (
	_ Stmt   = (*DeclStmt)(nil)
	_ Parent = (*DeclStmt)(nil)
)

func (n *DeclStmt) String() string { return nodeString(n.Decl) }

func (n *DeclStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*DeclStmt) StmtNode() {}

func (n *DeclStmt) Children(yield func(Node) bool) {
	_ = yield(n.Decl)
}
