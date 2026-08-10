package ir

import (
	"fmt"
	"go/token"
)

// RangeStmt is a node that represents a for statement with a range clause.
type RangeStmt struct {
	// ForPos is a position of "for" keyword
	ForPos token.Pos

	// Key is the first value from the range. May be nil.
	Key Expr

	// Value is the second value from the range. May be nil.
	Value Expr

	// Define is set true if the key and value are being created via this assignment
	// or set false if the variables exist and are being overwritten.
	// This has no affect if key and value.
	Define bool

	// X is the value being ranged over
	X Expr

	// Body is the statements in the body of the loop.
	Body []Stmt
}

var (
	_ Stmt   = (*RangeStmt)(nil)
	_ Parent = (*RangeStmt)(nil)
)

func (n *RangeStmt) String() string {
	def := `=`
	if n.Define {
		def = `:=`
	}
	head := ``
	if n.Key != nil {
		if n.Value != nil {
			head = fmt.Sprintf(`%s, %s %s `, toString(n.Key), toString(n.Value), def)
		} else {
			head = fmt.Sprintf(`%s %s `, toString(n.Key), def)
		}
	} else if n.Value != nil {
		head = fmt.Sprintf(`_, %s %s `, toString(n.Value), def)
	}
	return fmt.Sprintf(`for (%srange %s)%s`, head, toString(n.X), bodyString(n.Body))
}

func (n *RangeStmt) Pos() token.Pos { return n.ForPos }

func (*RangeStmt) StmtNode() {}

func (n *RangeStmt) Children(yield func(Node) bool) {
	_ = yield(n.Key) && yield(n.Value) && yield(n.X) && YieldSlice(n.Body, yield)
}
