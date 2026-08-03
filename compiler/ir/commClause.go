package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// CommClause is a node that represents a case of a select statement.
type CommClause struct {
	Ast  *ast.CommClause // TODO: REMOVE
	Comm Stmt            // send or receive statement; nil means default case
	Body []Stmt          // statement list; or nil
}

var (
	_ Node   = (*CommClause)(nil)
	_ Parent = (*CommClause)(nil)
)

func (n *CommClause) String() string {
	str := `default:`
	if n.Comm != nil {
		str = fmt.Sprintf(`case %v:`, n.Comm)
	}
	if len(n.Body) > 0 {
		str += "\n" + linesString(n.Body)
	}
	return str
}

func (n *CommClause) Pos() token.Pos { return astPos(n.Ast) }

func (n *CommClause) Children(yield func(Node) bool) {
	_ = yield(n.Comm) && YieldSlice(n.Body, yield)
}
