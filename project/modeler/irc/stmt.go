package irc

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

type (
	// Stmt is a code statement that can be put inside of a function
	// but has no type value like an expression.
	Stmt interface {
		// stmtNode is an empty method used to compile time type check that
		// only statements duck-type to this interface.
		stmtNode()
	}

	// BaseStmt is a statement containing Go AST statements.
	BaseStmt struct {
		Stmt ast.Stmt
	}
)

func (s *BaseStmt) String() string {
	buf := &bytes.Buffer{}
	fSet := token.NewFileSet()
	if err := format.Node(buf, fSet, s.Stmt); err != nil {
		panic(err)
	}
	return buf.String()
}

func (s *BaseStmt) stmtNode() {}
