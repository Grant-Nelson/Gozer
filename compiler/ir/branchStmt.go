package ir

import (
	"go/token"
	"go/types"

	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
)

// BranchStmt is a node that represents a break, continue, goto,
// or fallthrough statement.
type BranchStmt struct {
	// TokPos is the position of token for the branch kind.
	TokPos token.Pos

	// Tok is the kind of the branch.
	Kind branchKind.BranchKind

	// Label is the optional label to jump with or nil.
	Label *types.Label
}

var (
	_ Stmt     = (*BranchStmt)(nil)
	_ Ref      = (*BranchStmt)(nil)
	_ FlowCtrl = (*BranchStmt)(nil)
	_ Parent   = (*BranchStmt)(nil)
)

func (*BranchStmt) StmtNode()     {}
func (*BranchStmt) RefNode()      {}
func (*BranchStmt) FlowCtrlNode() {}

func (n *BranchStmt) Pos() token.Pos       { return n.TokPos }
func (n *BranchStmt) Type() types.Type     { return n.Label.Type() }
func (n *BranchStmt) Object() types.Object { return n.Label }

func (n *BranchStmt) String() string {
	if n.Label == nil {
		return toString(n.Kind)
	}
	return toString(n.Kind) + ` ` + n.Label.Name()
}

func (n *BranchStmt) Children(yield func(Node) bool) {
	_ = yield(n.Label)
}
