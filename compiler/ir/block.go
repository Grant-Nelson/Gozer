package ir

import (
	"fmt"
	"go/token"
	"strings"
)

// Block represents a statement block defining a set of statements.
// Typically a statement block should have no flow control (mainly loops)
// except for exiting the block or for inlined flow (like a small if-statement).
// See [README.md]
type Block struct {
	// Index is the ones-based index of this block in the full set of blocks
	// for a package. 0 or less if not assigned an index yet.
	Index int

	// Hint is a string to help debug block creation.
	Hint string

	// Body is the list of statements for this block.
	// A block is invalid when there are no statements.
	Body []Stmt

	// Prior blocks are blocks that can transition to this block.
	Prior []*Block

	// Follow blocks are blocks that this block can transition to.
	Follow []*Block

	// Params are the parameters that are passed into this block
	// when it is called and are available inside the block.
	//
	// This will not contain any parameters needed for the closure
	// including any receiver and any type dictionaries.
	Params []*Param
}

var (
	_ Node   = (*Block)(nil)
	_ Parent = (*Block)(nil)
)

func (n *Block) String() string {
	var params, hint, tail string
	if len(n.Params) > 0 {
		parts := make([]string, len(n.Params))
		for i, p := range n.Params {
			parts[i] = p.String()
		}
		params = strings.Join(parts, `, `)
	}
	if len(n.Hint) > 0 {
		hint = `<` + n.Hint + `> `
	}
	if len(n.Body) > 0 {
		tail = fmt.Sprintf("\n%s\n", linesString(n.Body))
	}
	return fmt.Sprintf("block %d (%s)%s{%s}", n.Index, params, hint, tail)
}

func (n *Block) Pos() token.Pos {
	if n == nil || len(n.Body) <= 0 || n.Body[0] == nil {
		return token.NoPos
	}
	return n.Body[0].Pos()
}

func (n *Block) Children(yield func(Node) bool) bool {
	return YieldSlice(n.Body, yield)
}

// LastStmt gets the last statement in the block or nil if empty.
func (n *Block) LastStmt() Stmt {
	if max := len(n.Body) - 1; max >= 0 {
		return n.Body[max]
	}
	return nil
}
