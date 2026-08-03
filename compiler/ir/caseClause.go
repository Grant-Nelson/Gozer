package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// CaseClause is a node that represents a case of an expression or type switch statement.
type CaseClause struct {
	Ast  *ast.CaseClause // TODO: REMOVE
	List []Expr          // list of expressions or types; nil means default case
	Body []Stmt          // statement list; or nil
}

var (
	_ Node   = (*CaseClause)(nil)
	_ Parent = (*CaseClause)(nil)
)

func (n *CaseClause) String() string {
	str := `default:`
	if len(n.List) > 0 {
		str = fmt.Sprintf(`case %v:`, csvString(n.List))
	}
	if len(n.Body) > 0 {
		str += "\n" + linesString(n.Body)
	}
	return str
}

func (n *CaseClause) Pos() token.Pos { return astPos(n.Ast) }

func (n *CaseClause) Children(yield func(Node) bool) {
	_ = YieldSlice(n.List, yield) && YieldSlice(n.Body, yield)
}
