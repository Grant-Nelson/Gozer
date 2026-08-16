package ir

import (
	"go/token"
	"go/types"
)

// FuncDecl represents a named function declaration.
type FuncDecl struct {

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
	Atomic bool

	FuncObj *types.Func

	// Func is the function definition.
	Func *Func
}

var (
	_ Stmt   = (*FuncDecl)(nil)
	_ Decl   = (*FuncDecl)(nil)
	_ Parent = (*FuncDecl)(nil)
)

func (*FuncDecl) StmtNode() {}
func (*FuncDecl) DeclNode() {}

func (fn *FuncDecl) Pos() token.Pos       { return fn.FuncObj.Pos() }
func (fn *FuncDecl) Type() types.Type     { return fn.Func.Signature }
func (fn *FuncDecl) Object() types.Object { return fn.FuncObj }

func (fn *FuncDecl) String() string {
	return `func ` + fn.FuncObj.FullName() + toString(fn.Func)
}

func (fn *FuncDecl) Children(yield func(Node) bool) {
	_ = yield(fn.Func)
}
