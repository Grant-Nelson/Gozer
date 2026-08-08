package ir

import (
	"fmt"
	"go/token"
)

// SwitchStmt is a node that represents an expression switch statement.
type SwitchStmt struct {
	SwitchPos token.Pos // position of "switch" keyword
	Init      Stmt      // initialization statement; or nil
	Tag       Expr      // tag expression; or nil
	Body      []*CaseClause
}

var _ Stmt = (*SwitchStmt)(nil)

func (n *SwitchStmt) String() string {
	return fmt.Sprintf("switch %s {\n%s\n}", n.Tag, linesString(n.Body))
}

func (n *SwitchStmt) Pos() token.Pos { return n.SwitchPos }

func (*SwitchStmt) StmtNode() {}

func (n *SwitchStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Tag) && YieldSlice(n.Body, yield)
}
