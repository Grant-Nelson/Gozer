package ir

import (
	"fmt"
	"go/token"
)

// CaseClause is a node that represents a case of an expression or type switch statement.
type CaseClause struct {

	// Case is the position of "case" or "default" keyword.
	Case token.Pos

	// List is the list of expressions or types; nil means default case
	List []Expr

	// Body is the statement list; or nil
	Body []Stmt
}

var (
	_ Stmt   = (*CaseClause)(nil)
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

func (*CaseClause) StmtNode() {}

func (n *CaseClause) Pos() token.Pos { return n.Case }

func (n *CaseClause) Children(yield func(Node) bool) {
	_ = YieldSlice(n.List, yield) && YieldSlice(n.Body, yield)
}

// IsDefault indicates if this is the default case for the switch.
func (n *CaseClause) IsDefault() bool { return len(n.List) == 0 }
