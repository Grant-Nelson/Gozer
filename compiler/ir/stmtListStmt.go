package ir

import (
	"go/token"
)

// StmtListStmt is a node that represents a braced statement list.
type StmtListStmt struct {
	List []Stmt
}

var (
	_ Stmt   = (*StmtListStmt)(nil)
	_ Parent = (*StmtListStmt)(nil)
)

func (*StmtListStmt) StmtNode() {}

func (n *StmtListStmt) Pos() token.Pos {
	if len(n.List) > 0 {
		return n.List[0].Pos()
	}
	return token.NoPos
}

func (n *StmtListStmt) String() string {
	return bodyString(n.List)
}

func (n *StmtListStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.List, yield)
}

func (n *StmtListStmt) Add(s Stmt) {
	if s != nil {
		n.List = append(n.List, s)
	}
}
