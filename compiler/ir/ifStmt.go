package ir

import (
	"go/token"
	"strings"
)

// An IfStmt node represents an if statement.
type IfStmt struct {
	IfPos token.Pos // position of "if" keyword
	Init  Stmt      // initialization statement; or nil
	Cond  Expr      // condition
	Body  []Stmt
	Else  []Stmt
}

var (
	_ Stmt   = (*IfStmt)(nil)
	_ Parent = (*IfStmt)(nil)
)

func (n *IfStmt) String() string {
	str := `if (` + toString(n.Cond) + `)` + bodyString(n.Body)
	if len(n.Else) <= 0 {
		return str
	}
	if strings.HasSuffix(str, `}`) {
		str += ` else`
	} else {
		str += "\nelse"
	}
	if len(n.Else) == 1 {
		if elseIf, ok := n.Else[0].(*IfStmt); ok {
			return str + ` ` + toString(elseIf)
		}
	}
	return str + bodyString(n.Else)
}

func (n *IfStmt) Pos() token.Pos { return n.IfPos }

func (*IfStmt) StmtNode() {}

func (n *IfStmt) Children(yield func(Node) bool) {
	_ = yield(n.Init) && yield(n.Cond) &&
		YieldSlice(n.Body, yield) && YieldSlice(n.Else, yield)
}
