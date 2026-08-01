package ir

import (
	"fmt"
	"go/ast"
	"go/token"
)

// RangeStmt is a node that represents a for statement with a range clause.
type RangeStmt struct {
	Ast   *ast.RangeStmt // TODO: REMOVE
	Key   ast.Expr       // TODO: REPLACE // Key may be nil
	Value ast.Expr       // TODO: REPLACE // Value may be nil
	Tok   token.Token    // TODO: REPLACE // ILLEGAL if Key == nil, ASSIGN, DEFINE
	X     ast.Expr       // TODO: REPLACE // value to range over
	Body  []Stmt
}

var _ Stmt = (*RangeStmt)(nil)

func (s *RangeStmt) String() string {
	return fmt.Sprintf("for %v, %v %v range %v {\n%v\n}",
		s.Key, s.Value, s.Tok.String(), s.X, linesString(s.Body))
}

func (s *RangeStmt) Pos() token.Pos { return astPos(s.Ast) }

func (*RangeStmt) StmtNode() {}

func FromRangeStmt(s *ast.RangeStmt, c *Converter) *RangeStmt {
	if s == nil {
		return nil
	}
	return &RangeStmt{
		Ast:   s,
		Key:   s.Key,
		Value: s.Value,
		X:     s.X,
		Body:  ExpandStmt(FromBlockStmt(s.Body, c)),
	}
}
