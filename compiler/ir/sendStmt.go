package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SendStmt is a node that represents a send statement.
type SendStmt struct {
	Ast   *ast.SendStmt // TODO: REMOVE
	Chan  ast.Expr      // TODO: REPLACE
	Value ast.Expr      // TODO: REPLACE
}

var _ Stmt = (*SendStmt)(nil)

func (s *SendStmt) String() string {
	return fmt.Sprintf(`%s<-%s`, nodeString(s.Chan), nodeString(s.Value))
}

func (s *SendStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*SendStmt) StmtNode() {}

func FromSendStmt(s *ast.SendStmt) *SendStmt {
	if s == nil {
		return nil
	}
	return &SendStmt{
		Ast:   s,
		Chan:  s.Chan,
		Value: s.Value,
	}
}
