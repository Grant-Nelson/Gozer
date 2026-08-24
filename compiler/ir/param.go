package ir

import "go/types"

// Parameter represents a single variable that can be passed into a block.
type Param struct {
	// Name is the identifier for this parameter.
	//
	// This identifier needs to have a types.Info entry to get the object.
	//
	// Unnamed parameters (parameter lists which only contain types) have
	// a nil name. This is not allowed here since these are only for blocks.
	// The function call to kick off the first block in a function will
	// only pass named parameters into the block. The function will keep
	// the unnamed parameters so that its signature remained the same as
	// it was defined in the AST.
	Name string

	// Type is the resolved type for this parameter.
	//
	// Used as a fallback when Expr is nil. The blocker sets this when
	// synthesizing block params for variables defined inside the function
	// body.
	Type types.Type
}

func (p *Param) String() string { return p.Name + ` ` + toString(p.Type) }
