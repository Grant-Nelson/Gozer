package ir

import (
	"fmt"
	"go/token"
)

// A LabeledStmt node represents a labeled statement.
type LabeledStmt struct {
	Label *Ident // same as AST so works with types.Info
	Stmt  Stmt   // TODO: INVERT RELATIONSHIP
}

var (
	_ Stmt   = (*LabeledStmt)(nil)
	_ Parent = (*LabeledStmt)(nil)
)

func (n *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", n.Label, n.Stmt) }

func (n *LabeledStmt) Pos() token.Pos { return n.Label.Pos() }

func (*LabeledStmt) StmtNode() {}

func (n *LabeledStmt) Children(yield func(Node) bool) {
	_ = yield(n.Label) && yield(n.Stmt)
}
