package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// BranchStmt is a node that represents a break, continue, goto,
// or fallthrough statement.
type BranchStmt struct {
	Ast   *ast.BranchStmt // TODO: REMOVE
	Tok   token.Token     // TODO: REPLACE // keyword token (BREAK, CONTINUE, GOTO, FALLTHROUGH)
	Label *ast.Ident      // TODO: REPLACE // label name; or nil
}

var (
	_ Stmt     = (*BranchStmt)(nil)
	_ FlowCtrl = (*BranchStmt)(nil)
	_ Parent   = (*BranchStmt)(nil)
)

func (n *BranchStmt) String() string { return fmt.Sprintf(`%s %v`, n.Tok.String(), n.Label) }

func (n *BranchStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*BranchStmt) StmtNode()     {}
func (*BranchStmt) FlowCtrlNode() {}

func (n *BranchStmt) Children(yield func(Node) bool) bool {
	return yield(n.Label)
}

func FromBranchStmt(s *ast.BranchStmt) *BranchStmt {
	if s == nil {
		return nil
	}
	return &BranchStmt{
		Ast:   s,
		Tok:   s.Tok,
		Label: s.Label,
	}
}
