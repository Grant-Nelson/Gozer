package ir

import (
	"go/token"
	"go/types"
)

// FuncLit represents a function literal that may be a closure.
//
// A function literal currently (based on Go specs) can not have a
// receiver nor type parameters. However, this literal may be defined
// within a closure that had type parameters so it may inherit
// type parameters and arguments.
type FuncLit struct {

	// Func is the function definition.
	Func *Func
}

var (
	_ Expr   = (*FuncLit)(nil)
	_ Parent = (*FuncLit)(nil)
)

func (fn *FuncLit) ExprNode() {}

func (fn *FuncLit) Pos() token.Pos   { return fn.Func.Pos() }
func (fn *FuncLit) Type() types.Type { return fn.Func.Signature }
func (fn *FuncLit) String() string   { return `funcLit ` + toString(fn.Func) }

func (fn *FuncLit) Children(yield func(Node) bool) { _ = yield(fn.Func) }
