package ir

import (
	"fmt"
	"go/token"
)

// RangeStmt is a node that represents a for statement with a range clause.
type RangeStmt struct {
	ForPos token.Pos   // position of "for" keyword
	Key    Expr        // Key may be nil
	Value  Expr        // Value may be nil
	Tok    token.Token // TODO: REPLACE // ILLEGAL if Key == nil, ASSIGN, DEFINE
	X      Expr        // value to range over
	Body   []Stmt
}

var (
	_ Stmt   = (*RangeStmt)(nil)
	_ Parent = (*RangeStmt)(nil)
)

func (n *RangeStmt) String() string {
	return fmt.Sprintf("for %v, %v %v range %v {\n%v\n}",
		n.Key, n.Value, n.Tok.String(), n.X, linesString(n.Body))
}

func (n *RangeStmt) Pos() token.Pos { return n.ForPos }

func (*RangeStmt) StmtNode() {}

func (n *RangeStmt) Children(yield func(Node) bool) {
	_ = yield(n.Key) && yield(n.Value) && yield(n.X) && YieldSlice(n.Body, yield)
}
