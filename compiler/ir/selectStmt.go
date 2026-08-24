package ir

import "go/token"

// SelectStmt is a node that represents a select statement.
type SelectStmt struct {
	SelectPos token.Pos // position of "select" keyword
	Body      []*CommClause
}

var (
	_ Stmt   = (*SelectStmt)(nil)
	_ Parent = (*SelectStmt)(nil)
)

func (*SelectStmt) StmtNode() {}

func (n *SelectStmt) Pos() token.Pos { return n.SelectPos }

func (n *SelectStmt) String() string {
	return `select {` + nlStr + linesString(n.Body) + nlStr + `}`
}

func (n *SelectStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Body, yield)
}
