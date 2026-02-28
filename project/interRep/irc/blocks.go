package irc

import (
	"fmt"
	"go/ast"
	"strings"
)

type (
	Package struct {
		Funcs []*Func
	}

	// Func represents a function block defining a function as
	// a collection of statement blocks.
	// See [README.md]
	Func struct {
		// Ast is the original AST for this function.
		Ast *ast.FuncDecl

		// Blocks is the collection of statement blocks for this function.
		// The first block in this slice is the entry point for this function.
		Blocks []*Block
	}

	// Block represents a statement block defining a set of statements.
	// Typically a statement block should have no flow control (mainly loops)
	// except for exiting the block or for inlined flow (like a small if-statement).
	// See [README.md]
	Block struct {
		// Index is the ones-based index of this block in the full set of blocks
		// for a package. 0 or less if not assigned an index yet.
		Index int

		// Prior blocks are blocks that can transition to this block.
		Prior []*Block

		// Follow blocks are blocks that this block can transition to.
		Follow []*Block

		// ExitsFunc indicates if this block can exit the function.
		ExitsFunc bool
	}

	// BlockRef is a reference for a block invocation.
	// See [docs/Blocks.md] for more information.
	BlockRef struct {
		Block *Block // block being referenced for invocation
		Args  []Expr // the arguments to pass onto the block when invoked
	}
)

func (b *Block) String() string {
	// TODO: Finish
	return `BLOCK`
}

func (s *BlockRef) String() string { return fmt.Sprintf(`%s, [%s]`, s.Block, csvString(s.Args)) }

func csvString[E any, S []E](s S) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		elems[i] = fmt.Sprintf(`%v`, elem)
	}
	return strings.Join(elems, `, `)
}

func linesString[E any, S []E](s S, indent string) string {
	elems := make([]string, len(s))
	for i, elem := range s {
		eStr := fmt.Sprintf(`%s%v`, indent, elem)
		elems[i] = strings.ReplaceAll(eStr, "\n", "\n"+indent)
	}
	return strings.Join(elems, "\n")
}
