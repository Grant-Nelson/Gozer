package ir

import (
	"go/ast"
	"go/types"
	"reflect"

	"github.com/Grant-Nelson/Gozer/avail/faults"
)

type Converter struct {
	Info *types.Info
}

func isNotNil[T any](t T) bool {
	v := reflect.ValueOf(t)
	return v.IsValid() && !v.IsZero()
}

func fromNilSafeStmt[TIn ast.Stmt, TOut Stmt](s TIn, fn func(TIn) TOut) Stmt {
	if isNotNil(s) {
		if t := fn(s); isNotNil(t) {
			return t
		}
	}
	return nil
}

func fromNilSafeStmt2[TIn1 ast.Stmt, TIn2 any, TOut Stmt](s TIn1, c TIn2, fn func(TIn1, TIn2) TOut) Stmt {
	if isNotNil(s) {
		if t := fn(s, c); isNotNil(t) {
			return t
		}
	}
	return nil
}

func FromStmt(s ast.Stmt, c *Converter) Stmt {
	switch s := s.(type) {
	case nil, *ast.BadStmt, *ast.EmptyStmt:
		return nil
	case *ast.DeclStmt:
		return fromNilSafeStmt(s, FromDeclStmt)
	case *ast.LabeledStmt:
		return fromNilSafeStmt2(s, c, FromLabeledStmt)
	case *ast.ExprStmt:
		return fromNilSafeStmt(s, FromExprStmt)
	case *ast.SendStmt:
		return fromNilSafeStmt(s, FromSendStmt)
	case *ast.IncDecStmt:
		return fromNilSafeStmt2(s, c, FromIncDecStmt)
	case *ast.AssignStmt:
		return fromNilSafeStmt(s, FromAssignStmt)
	case *ast.GoStmt:
		return fromNilSafeStmt(s, FromGoStmt)
	case *ast.DeferStmt:
		return fromNilSafeStmt(s, FromDeferStmt)
	case *ast.ReturnStmt:
		return fromNilSafeStmt(s, FromReturnStmt)
	case *ast.BranchStmt:
		return fromNilSafeStmt(s, FromBranchStmt)
	case *ast.BlockStmt:
		return fromNilSafeStmt2(s, c, FromBlockStmt)
	case *ast.IfStmt:
		return fromNilSafeStmt2(s, c, FromIfStmt)
	case *ast.SwitchStmt:
		return fromNilSafeStmt2(s, c, FromSwitchStmt)
	case *ast.TypeSwitchStmt:
		return fromNilSafeStmt2(s, c, FromTypeSwitchStmt)
	case *ast.SelectStmt:
		return fromNilSafeStmt2(s, c, FromSelectStmt)
	case *ast.ForStmt:
		return fromNilSafeStmt2(s, c, FromForStmt)
	case *ast.RangeStmt:
		return fromNilSafeStmt2(s, c, FromRangeStmt)
	default:
		panic(faults.New(`unexpected AST statement type`).
			WithF(`type`, `%T`, s))
	}
}
