package ir

import (
	"fmt"
	"go/token"
)

// CommClause is a node that represents a case of a select statement.
type CommClause struct {

	// Case is the position of "case" or "default" keyword.
	Case token.Pos

	// Comm is the send or receive statement; nil means default case
	Comm Stmt

	// Body is the statement list; or nil
	Body []Stmt
}

var (
	_ Stmt   = (*CommClause)(nil)
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

func (*CommClause) StmtNode() {}

func (n *CommClause) Pos() token.Pos { return n.Case }

func (n *CommClause) Children(yield func(Node) bool) {
	_ = yield(n.Comm) && YieldSlice(n.Body, yield)
}

// IsDefault indicates if this is the default case for the switch.
func (n *CommClause) IsDefault() bool { return n.Comm == nil }
