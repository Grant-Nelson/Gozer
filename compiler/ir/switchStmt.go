package ir

import "go/token"

// SwitchStmt is a node that represents an expression switch statement.
type SwitchStmt struct {
	SwitchPos token.Pos // position of "switch" keyword
	Init      Stmt      // initialization statement; or nil
	Tag       Expr      // tag expression; or nil
	Body      []*CaseClause
}

var (
	_ Stmt   = (*SwitchStmt)(nil)
	_ Parent = (*SwitchStmt)(nil)
)

func (*SwitchStmt) StmtNode() {}

func (n *SwitchStmt) Pos() token.Pos { return n.SwitchPos }

func (n *SwitchStmt) String() string {
	return `switch ` + toString(n.Tag) + bodyString(n.Body)
}

func (n *SwitchStmt) ChildCount() int { return 2 + len(n.Body) }

func (n *SwitchStmt) Children(yield func(Node) bool) {
	_ = YieldNode(n.Init, yield) &&
		YieldNode(n.Tag, yield) &&
		YieldSlice(n.Body, yield)
}
