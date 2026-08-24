package ir

import (
	"go/token"
	"go/types"
)

// A LabeledStmt node represents a labeled statement.
//
// A label acts like a declaration since it declares an object
// that is referenced by break, continue, and goto statements.
//
// Labels becomes part of a block with the statement defining how the
// block looks. A goto or continue will jump to the statement's block
// and a break will jump to after the statement. If there is no statement
// then the label indicates the start of a block that a goto can jump to.
type LabeledStmt struct {
	LabelObj *types.Label
	Stmt     Stmt
}

var (
	_ Stmt   = (*LabeledStmt)(nil)
	_ Decl   = (*LabeledStmt)(nil)
	_ Parent = (*LabeledStmt)(nil)
)

func (*LabeledStmt) StmtNode() {}
func (*LabeledStmt) DeclNode() {}

func (n *LabeledStmt) Pos() token.Pos       { return n.LabelObj.Pos() }
func (n *LabeledStmt) Type() types.Type     { return n.LabelObj.Type() }
func (n *LabeledStmt) Object() types.Object { return n.LabelObj }

func (n *LabeledStmt) String() string {
	if n.Stmt == nil {
		return n.LabelObj.Name() + `:`
	}
	return n.LabelObj.Name() + `:` + nlStr + toString(n.Stmt)
}

func (n *LabeledStmt) Children(yield func(Node) bool) { _ = yield(n.Stmt) }
