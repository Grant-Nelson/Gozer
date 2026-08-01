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

var (
	_ Stmt   = (*SendStmt)(nil)
	_ Parent = (*SendStmt)(nil)
)

func (n *SendStmt) String() string {
	return fmt.Sprintf(`%s<-%s`, nodeString(n.Chan), nodeString(n.Value))
}

func (n *SendStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*SendStmt) StmtNode() {}

func (n *SendStmt) Children(yield func(Node) bool) bool {
	return yield(n.Chan) && yield(n.Value)
}

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
