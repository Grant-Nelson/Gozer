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

var (
	_ Stmt   = (*RangeStmt)(nil)
	_ Parent = (*RangeStmt)(nil)
)

func (n *RangeStmt) String() string {
	return fmt.Sprintf("for %v, %v %v range %v {\n%v\n}",
		n.Key, n.Value, n.Tok.String(), n.X, linesString(n.Body))
}

func (n *RangeStmt) Pos() token.Pos { return astPos(n.Ast) }

func (*RangeStmt) StmtNode() {}

func (n *RangeStmt) Children(yield func(Node) bool) bool {
	return yield(n.Key) && yield(n.Value) && yield(n.X) && YieldSlice(n.Body, yield)
}

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
