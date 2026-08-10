package ir

import (
	"fmt"
	"go/token"
)

// MultiAssignStmt is a node that represents an assignment or
// a short variable declaration.
type MultiAssignStmt struct {

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
	_ Stmt   = (*MultiAssignStmt)(nil)
	_ Parent = (*MultiAssignStmt)(nil)
)

func (n *MultiAssignStmt) String() string {
	def := ` = `
	if n.Define {
		def = ` := `
	}
	return fmt.Sprintf(`%s%s%s`, csvString(n.Lhs), def, csvString(n.Rhs))
}

func (n *MultiAssignStmt) Pos() token.Pos { return n.TokPos }

func (*MultiAssignStmt) StmtNode() {}

func (n *MultiAssignStmt) Children(yield func(Node) bool) {
	_ = YieldSlice(n.Lhs, yield) && YieldSlice(n.Rhs, yield)
}
