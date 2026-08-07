package ir

import (
	"fmt"
	"go/token"
)

// InjectStmt is a statement containing code in the target language that
// can be used to drop explicitly defined target code into the output.
type InjectStmt struct {
	// Source is an optional position for where this data is coming from.
	Source token.Pos

	// Injection is the code in the target language to drop in at this location.
	Injection string
}

var _ Stmt = (*InjectStmt)(nil)

func (n *InjectStmt) String() string {
	return fmt.Sprintf(`inject[%q]`, n.Injection)
}

func (n *InjectStmt) Pos() token.Pos { return n.Source }

func (*InjectStmt) StmtNode() {}
