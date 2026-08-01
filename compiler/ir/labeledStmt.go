package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// A LabeledStmt node represents a labeled statement.
type LabeledStmt struct {
	Ast   *ast.LabeledStmt // TODO: REMOVE
	Label *ast.Ident       // TODO: REPLACE // same as AST so works with types.Info
	Stmt  Stmt             // TODO: INVERT RELATIONSHIP
}

var (
	_ Stmt   = (*LabeledStmt)(nil)
	_ Parent = (*LabeledStmt)(nil)
)

func (n *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", n.Label, n.Stmt) }

func (n *LabeledStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*LabeledStmt) StmtNode() {}

func (n *LabeledStmt) Children(yield func(Node) bool) bool {
	return yield(n.Label) && yield(n.Stmt)
}

func FromLabeledStmt(s *ast.LabeledStmt, c *Converter) *LabeledStmt {
	if s == nil {
		return nil
	}
	return &LabeledStmt{
		Ast:   s,
		Label: s.Label,
		Stmt:  FromStmt(s.Stmt, c),
	}
}
