package ir

import (
	"go/token"
	"go/types"
)

// Func represents a function block defining a function as
// a collection of statement blocks.
type Func struct {

	// FuncPos is the position of `func` keyword for the function, method,
	// or function literal.
	FuncPos token.Pos

	// Signature is the type for the signature of the function.
	Signature *types.Signature

	// Blocks is the collection of statement blocks for this function.
	// The first block in this slice is the entry point for this function.
	Blocks []*Block

	// ReturnBlocks are blocks that exit this function.
	ReturnBlocks []*Block
}

var (
	_ Parent = (*Func)(nil)
	_ Parent = (*Func)(nil)
)

func (fn *Func) Pos() token.Pos { return fn.FuncPos }

func (fn *Func) String() string {
	if len(fn.Blocks) <= 0 {
		return `{ }`
	}
	return `{` + nlStr + linesString(fn.Blocks) + nlStr + `}`
}

func (fn *Func) Children(yield func(Node) bool) {
	_ = YieldSlice(fn.Blocks, yield)
}

// NewBlock creates a new empty block and adds it to this function.
func (fn *Func) NewBlock(hint string, body []Stmt, params []*Param) *Block {
	b := &Block{
		Index:  len(fn.Blocks),
		Hint:   hint,
		Body:   body,
		Params: params,
	}
	fn.Blocks = append(fn.Blocks, b)
	return b
}
