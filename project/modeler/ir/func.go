package ir

import (
	"fmt"
	"go/ast"

	"github.com/Grant-Nelson/Gozer/avail/assert"
	"github.com/Grant-Nelson/Gozer/avail/astTools"
)

// Func represents a function block defining a function as
// a collection of statement blocks.
// See [README.md]
type Func struct {
	// Package is the package this function belongs to.
	Package *Package

	// Ast is the AST for this function.
	//
	// Remodelers that modify the AST inside of a [BaseStmt] or similar
	// may cause changes to this data as well.
	Ast *ast.FuncDecl

	// TODO: Add a signature with any receivers as the first parameters.

	// Name of the function initialized to the name from the AST.
	//
	// The function name is unique per package. If the function has a receiver,
	// then that receiver type is used to prefix the name.
	// If this name collides with another top-level object,
	// the name maybe modified to be unique, or modified to fit the
	// target language style better.
	Name string

	// Blocks is the collection of statement blocks for this function.
	// The first block in this slice is the entry point for this function.
	Blocks []*Block
}

func (fn *Func) String() string {
	return fmt.Sprintf("func %s {\n%s\n}", fn.Name, linesString(fn.Blocks))
}

// NewBlock creates a new empty block and adds it to this function.
func (fn *Func) NewBlock() *Block {
	b := &Block{Index: len(fn.Blocks)}
	fn.Blocks = append(fn.Blocks, b)
	return b
}

// Atomic indicates that this function should NOT be broken up into
// flow control blocks, meaning that the scheduler running in one real
// thread environment should not swap goroutines.
//
// Atomic is based on the directive `//gozer:atomic`.
// However, the function may still not be atomic if it contains blocking
// calls such as a send, receive, sleep, or lock. Calling a function
// on an interface that is not pinned to a package could be blocking so
// is treated as always blocking. That goes for calling a function on
// a generic type as well.
//
// For a target language like typescript, calling a non-atomic function
// from an atomic function will have a signature mismatch since the parameters
// and returns from a non-atomic function will have parameters and returns
// specifically designed for the schedular to call.
func (fn *Func) Atomic() bool {
	if fn.Ast.Doc == nil {
		return false
	}
	dv := astTools.Directives(fn.Ast.Doc.List, directiveGroup)
	if s, ok := dv[directiveAtomicFunc]; ok {
		// The atomic directive should have no fields.
		// Any fields will be ignored (if asserts are off).
		assert.EmptySlice(s)
		return true
	}
	return false
}
