package ir

import "go/token"

// Node represents any part of the IR Tree.
type Node interface {

	// Pos is the position information for the source location
	// this node came from.
	Pos() token.Pos

	// String returns a human-readable text representing the node
	// for debugging and testing. The output must be consistent.
	String() string
}
