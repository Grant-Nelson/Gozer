package ir

// Stmt is a code statement that can be put inside of a function
// but has no type value like an expression.
type Stmt interface {
	Node

	// StmtNode is an empty method used for duck-typing statements.
	StmtNode()
}

var _ Node = (Stmt)(nil)
