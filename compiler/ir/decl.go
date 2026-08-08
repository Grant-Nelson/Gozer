package ir

// Decl is a node for a declaration.
type Decl interface {
	Node

	// DeclNode is an empty method used to compile time type check
	// that only declarations duck-type to this interface.
	DeclNode()
}
