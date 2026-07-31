package ir

type (
	// Stmt is a code statement that can be put inside of a function
	// but has no type value like an expression.
	Expr interface {
		Node

		// String returns a human-readable text representing the statement
		// for debugging and testing. The output must be consistent.
		String() string

		// ExprNode is an empty method used to compile time type check that
		// only expression duck-type to this interface.
		ExprNode()
	}
)
