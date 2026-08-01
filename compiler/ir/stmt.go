package ir

import "go/ast"

// Stmt is a code statement that can be put inside of a function
// but has no type value like an expression.
type Stmt interface {
	Node

	// StmtNode is an empty method used for duck-typing statements.
	StmtNode()
}

var _ Node = (Stmt)(nil)

func FromStmtSlice(ss []ast.Stmt, c *Converter) []Stmt {
	result := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		result = append(result, ExpandStmt(FromStmt(s, c))...)
	}
	return result
}

func ExpandStmtSlice(ss []Stmt) []Stmt {
	st := make([]Stmt, 0, len(ss))
	for _, s := range ss {
		st = append(st, ExpandStmt(s)...)
	}
	return st
}

func ExpandStmt(s Stmt) []Stmt {
	if s == nil {
		return []Stmt{}
	}
	if b, ok := s.(*StmtListStmt); ok {
		return ExpandStmtSlice(b.List)
	}
	return []Stmt{s}
}
