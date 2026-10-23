package ir

import "go/token"

// InjectStmt is a statement containing code in the target language that
// can be used to drop explicitly defined target code into the output.
type InjectStmt struct {
	// Source is an optional position for where this data is coming from.
	Source token.Pos

	// Injection is the code in the target language to drop in at this location.
	Injection string
}

var _ Stmt = (*InjectStmt)(nil)

func (*InjectStmt) StmtNode() {}

func (n *InjectStmt) Pos() token.Pos { return n.Source }

func (n *InjectStmt) String() string {
	return `inject[` + toString(n.Injection) + `]`
}
