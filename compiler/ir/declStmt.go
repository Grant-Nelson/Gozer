package ir

import (
	"go/ast"
	"go/token"
)

// DeclStmt is a node that represents a declaration in a statement list.
type DeclStmt struct {
	Ast  *ast.DeclStmt // TODO: REMOVE
	Decl ast.Decl      // TODO: REPLACE // *GenDecl with CONST, TYPE, or VAR token
}

var _ Stmt = (*DeclStmt)(nil)

func (s *DeclStmt) String() string { return nodeString(s.Decl) }

func (s *DeclStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*DeclStmt) StmtNode() {}

func FromDeclStmt(s *ast.DeclStmt) *DeclStmt {
	if s == nil {
		return nil
	}
	return &DeclStmt{
		Ast:  s,
		Decl: s.Decl,
	}
}
