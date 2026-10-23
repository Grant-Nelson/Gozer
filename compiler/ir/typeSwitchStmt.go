package ir

import "go/token"

// TypeSwitchStmt is a node that represents a type switch statement.
type TypeSwitchStmt struct {
	SwitchPos token.Pos // position of "switch" keyword
	Init      Stmt      // initialization statement; or nil
	Assign    Stmt      // x := y.(type) or y.(type)
	Body      []*CaseClause
}

var (
	_ Stmt   = (*TypeSwitchStmt)(nil)
	_ Parent = (*TypeSwitchStmt)(nil)
)

func (*TypeSwitchStmt) StmtNode() {}

func (n *TypeSwitchStmt) Pos() token.Pos { return n.SwitchPos }

func (n *TypeSwitchStmt) String() string {
	return `switch ` + toString(n.Assign) + bodyString(n.Body)
}

func (n *TypeSwitchStmt) ChildCount() int { return 2 + len(n.Body) }

func (n *TypeSwitchStmt) Children(yield func(Node) bool) {
	_ = YieldNode(n.Init, yield) &&
		YieldNode(n.Assign, yield) &&
		YieldSlice(n.Body, yield)
}
