package ir

import (
	"fmt"
	"go/token"
)

// SendStmt is a node that represents a send statement.
type SendStmt struct {
	ArrowPos token.Pos // position of "<-"
	Chan     Expr
	Value    Expr
}

var (
	_ Stmt   = (*SendStmt)(nil)
	_ Parent = (*SendStmt)(nil)
)

func (n *SendStmt) String() string {
	return fmt.Sprintf(`%s<-%s`, n.Chan, n.Value)
}

func (n *SendStmt) Pos() token.Pos { return n.ArrowPos }

func (*SendStmt) StmtNode() {}

func (n *SendStmt) Children(yield func(Node) bool) {
	_ = yield(n.Chan) && yield(n.Value)
}
