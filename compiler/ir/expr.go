package ir

// Expr is a code statement that can be put inside of a
// function but has no type value like an expression.
type Expr interface {
	Node

	// ExprNode is an empty method used to compile time type check
	// that only expression duck-type to this interface.
	ExprNode()
}
