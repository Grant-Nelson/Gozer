package ir

import "go/types"

// Expr is a code statement that can be put inside of a
// function but has no type value like an expression.
//
// The IR allows all expressions to be used as statements
// to allow for target language constructs that Go restricts.
// This allows for statements like `x = y = z++`.
type Expr interface {
	Stmt

	// Type is the resulting type of this expression.
	// If this returns nil, then there is no resulting type.
	Type() types.Type

	// ExprNode is an empty method used to compile time type check
	// that only expression duck-type to this interface.
	ExprNode()
}
