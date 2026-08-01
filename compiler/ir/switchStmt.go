package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// SwitchStmt is a node that represents an expression switch statement.
type SwitchStmt struct {
	Ast  *ast.SwitchStmt
	Init Stmt     // initialization statement; or nil
	Tag  ast.Expr // tag expression; or nil
	Body []*CaseClause
}

var _ Stmt = (*SwitchStmt)(nil)

func (n *SwitchStmt) String() string {
	return fmt.Sprintf("switch %s {\n%s\n}", nodeString(n.Tag), linesString(n.Body))
}

func (n *SwitchStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*SwitchStmt) StmtNode() {}

func (n *SwitchStmt) Children(yield func(Node) bool) bool {
	return yield(n.Init) && yield(n.Tag) && YieldSlice(n.Body, yield)
}

func FromSwitchStmt(s *ast.SwitchStmt, c *Converter) *SwitchStmt {
	if s == nil {
		return nil
	}
	return &SwitchStmt{
		Ast:  s,
		Init: FromStmt(s.Init, c),
		Tag:  s.Tag,
		Body: FromCaseClauseSlice(s.Body, c),
	}
}
