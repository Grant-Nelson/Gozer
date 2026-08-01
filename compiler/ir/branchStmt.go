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

var _ Stmt = (*BranchStmt)(nil)

func (s *BranchStmt) String() string { return fmt.Sprintf(`%s %v`, s.Tok.String(), s.Label) }

func (s *BranchStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*BranchStmt) StmtNode() {}

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
