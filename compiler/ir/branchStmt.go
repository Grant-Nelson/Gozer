package ir

import (
	"fmt"
	"go/ast"
	"go/token"

	"github.com/Grant-Nelson/Gozer/compiler/ir/enums/branchKind"
)

// BranchStmt is a node that represents a break, continue, goto,
// or fallthrough statement.
type BranchStmt struct {
	TokPos token.Pos             // TokPos is the position of token for the branch kind.
	Kind   branchKind.BranchKind // Tok is the kind of the branch.
	Label  *ast.Ident            // TODO: REPLACE // label name; or nil
}

var (
	_ Stmt     = (*BranchStmt)(nil)
	_ FlowCtrl = (*BranchStmt)(nil)
	_ Parent   = (*BranchStmt)(nil)
)

func (n *BranchStmt) String() string { return fmt.Sprintf(`%s %v`, n.Kind.String(), n.Label) }

func (n *BranchStmt) Pos() token.Pos { return n.TokPos }

func (*BranchStmt) StmtNode()     {}
func (*BranchStmt) FlowCtrlNode() {}

func (n *BranchStmt) Children(yield func(Node) bool) {
	_ = yield(n.Label)
}
