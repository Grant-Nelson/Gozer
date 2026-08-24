package ir

import "go/token"

// TypeSwitchStmt is a node that represents a type switch statement.
type TypeSwitchStmt struct {
	SwitchPos token.Pos // position of "switch" keyword
	Init      Stmt      // initialization statement; or nil
	Assign    Stmt      // x := y.(type) or y.(type)
	Body      []*CaseClause
}

var _ Stmt = (*TypeSwitchStmt)(nil)

func (*TypeSwitchStmt) StmtNode() {}

func (n *TypeSwitchStmt) Pos() token.Pos { return n.SwitchPos }

func (n *TypeSwitchStmt) String() string {
	return `switch ` + toString(n.Assign) + ` {` + nlStr + linesString(n.Body) + nlStr + `}`
}

func (n *TypeSwitchStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Assign) && YieldSlice(n.Body, yield)
}
