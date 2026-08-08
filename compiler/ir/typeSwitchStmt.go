package ir

import (
	"fmt"
	"go/token"
)

// TypeSwitchStmt is a node that represents a type switch statement.
type TypeSwitchStmt struct {
	SwitchPos token.Pos // position of "switch" keyword
	Init      Stmt      // initialization statement; or nil
	Assign    Stmt      // x := y.(type) or y.(type)
	Body      []*CaseClause
}

var _ Stmt = (*TypeSwitchStmt)(nil)

func (n *TypeSwitchStmt) String() string {
	return fmt.Sprintf("switch %v {\n%s\n}", n.Assign, linesString(n.Body))
}

func (n *TypeSwitchStmt) Pos() token.Pos { return n.SwitchPos }

func (*TypeSwitchStmt) StmtNode() {}

func (n *TypeSwitchStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Assign) && YieldSlice(n.Body, yield)
}
