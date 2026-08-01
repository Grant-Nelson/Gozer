package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// ReturnStmt is a node that represents a return statement.
type ReturnStmt struct {
	Ast     *ast.ReturnStmt // TODO: REMOVE
	Results []ast.Expr      // TODO: REPLACE // result expressions; or nil
}

var (
	_ Stmt     = (*ReturnStmt)(nil)
	_ FlowCtrl = (*ReturnStmt)(nil)
)

func (s *ReturnStmt) String() string {
	if len(s.Results) <= 0 {
		return `return`
	}
	return fmt.Sprintf(`return %s`, csvString(s.Results))
}

func (s *ReturnStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*ReturnStmt) StmtNode()     {}
func (*ReturnStmt) FlowCtrlNode() {}

func FromReturnStmt(s *ast.ReturnStmt) *ReturnStmt {
	if s == nil {
		return nil
	}
	return &ReturnStmt{
		Ast:     s,
		Results: s.Results,
	}
}
