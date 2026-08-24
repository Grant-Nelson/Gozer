package ir

import "go/token"

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

func (*SendStmt) StmtNode() {}

func (n *SendStmt) Pos() token.Pos { return n.ArrowPos }

func (n *SendStmt) String() string {
	return toString(n.Chan) + `<-` + toString(n.Value)
}

func (n *SendStmt) Children(yield func(Node) bool) {
	_ = yield(n.Chan) && yield(n.Value)
}
