package ir

import (
	"fmt"
	"go/token"
)

// TODO: Update to MultiAssign and use BinaryExpr for single expression and definitions

// AssignStmt is a node that represents an assignment or
// a short variable declaration.
type AssignStmt struct {

	// TokPos is the position for the assignment token, `=` or `:=`.
	TokPos token.Pos

	// Lhs are the left hand side of the assignment for the
	// variables being assigned or defined.
	Lhs []Expr

	// Define is set true if the variables are being created via this assignment
	// or set false if the variables exist and are being overwritten.
	Define bool

	// Rhs are the right hand side for the values to assign to the left hand side.
	Rhs []Expr
}

var (
	_ Stmt   = (*AssignStmt)(nil)
	_ Parent = (*AssignStmt)(nil)
)

func (n *AssignStmt) String() string {
	def := ` = `
	if n.Define {
		def = ` := `
	}
	return fmt.Sprintf(`%s%s%s`, csvString(n.Lhs), def, csvString(n.Rhs))
}

func (n *AssignStmt) Pos() token.Pos { return n.TokPos }

func (*AssignStmt) StmtNode() {}

func (n *AssignStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Lhs, yield) && YieldSlice(n.Rhs, yield)
}
