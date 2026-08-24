package ir

import "go/token"

// ReturnStmt is a node that represents a return statement.
type ReturnStmt struct {
	ReturnPos token.Pos // position of "return" keyword
	Results   []Expr    // result expressions; or nil
}

var (
	_ Stmt     = (*ReturnStmt)(nil)
	_ FlowCtrl = (*ReturnStmt)(nil)
	_ Parent   = (*ReturnStmt)(nil)
)

func (*ReturnStmt) StmtNode()     {}
func (*ReturnStmt) FlowCtrlNode() {}

func (n *ReturnStmt) Pos() token.Pos { return n.ReturnPos }

func (n *ReturnStmt) String() string {
	if len(n.Results) <= 0 {
		return `return`
	}
	return `return ` + csvString(n.Results)
}

func (n *ReturnStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Results, yield)
}
