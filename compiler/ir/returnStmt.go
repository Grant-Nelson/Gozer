package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ReturnStmt is a node that represents a return statement.
type ReturnStmt struct {
	Ast     *ast.ReturnStmt // TODO: REMOVE
	Results []Expr          // result expressions; or nil
}

var (
	_ Stmt     = (*ReturnStmt)(nil)
	_ FlowCtrl = (*ReturnStmt)(nil)
	_ Parent   = (*ReturnStmt)(nil)
)

func (n *ReturnStmt) String() string {
	if len(n.Results) <= 0 {
		return `return`
	}
	return fmt.Sprintf(`return %s`, csvString(n.Results))
}

func (n *ReturnStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*ReturnStmt) StmtNode()     {}
func (*ReturnStmt) FlowCtrlNode() {}

func (n *ReturnStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Results, yield)
}
