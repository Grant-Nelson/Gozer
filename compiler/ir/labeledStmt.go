package ir

import (
	"fmt"
	"go/token"
)

// A LabeledStmt node represents a labeled statement.
type LabeledStmt struct {
	Name    string
	NamePos token.Pos
	Stmt    Stmt
}

var (
	_ Stmt   = (*LabeledStmt)(nil)
	_ Parent = (*LabeledStmt)(nil)
)

func (n *LabeledStmt) String() string { return fmt.Sprintf("%s:\n%v", n.Name, n.Stmt) }

func (n *LabeledStmt) Pos() token.Pos { return n.NamePos }

func (*LabeledStmt) StmtNode() {}

func (n *LabeledStmt) Children(yield func(Node) bool) { _ = yield(n.Stmt) }
