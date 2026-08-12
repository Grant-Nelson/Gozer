package ir

import (
	"go/token"
	"go/types"
)

// FuncLit represents a function literal.
type FuncLit struct {
	Func *Func

	FuncType types.Type
}

var _ Expr = (*FuncLit)(nil)

func (fn *FuncLit) Pos() token.Pos { return fn.Func.Pos() }

func (fn *FuncLit) Type() types.Type { return fn.FuncType }

func (fn *FuncLit) ExprNode() {}

func (fn *FuncLit) String() string { return `funcLit: ` + toString(fn.Func) }
